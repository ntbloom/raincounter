/* arduino_rain_gauge.ino
 *
 * Increment rain gauge counter on click, display on screen
 *
 */


#include <Arduino.h>

#include "src/button.hpp"
#include "src/raingauge.hpp"
#include "src/serial/static_serial_tlv.hpp"
#include "src/temp36.hpp"
#include "src/timer.hpp"

#include "src/boards/mkr1000.cpp"
#include "src/config.cpp"

/* vendoring main from ArduinoCore-samd */
void initVariant() __attribute__((weak));
void initVariant()
{
}
extern USBDeviceClass USBDevice;
extern "C" void __libc_init_array(void);

using components::Button;
using components::Raingauge;
using components::Temp36;
using tlv::StaticSerialTLV;
using utilities::Timer;

/* define variables */
bool paused = false;

/* initialize components */
Button *holdButton = new Button(PAUSE_PIN, 50, HIGH);
Temp36 *tempSensor = new Temp36(TEMP_PIN, TEMP_VOLTAGE);
Timer *tempTimer = new Timer(TEMP_INTERVAL);
Raingauge *rainGauge = new Raingauge(RAIN_PIN, 50, GAUGE_MET, GAUGE_STD);
StaticSerialTLV *serialTLV = new StaticSerialTLV();

/* increment the rain counters */
void handleRainGauge()
{
    rainGauge->addCount();
    serialTLV->sendRainEvent();
}

/* don't increment counters on click, display pause message on screen */
void handlePause()
{
    if (paused)
    {
        // unpause the screen
        paused = false;
        serialTLV->sendUnpause();
        digitalWrite(LED_RED, 0);
    }
    else
    {
        paused = true;
        serialTLV->sendPause();
        digitalWrite(LED_RED, 1);
    }
}

/* take temperature measurement */
void handleMeasureTemp()
{
    digitalWrite(LED_BLUE, 1);
    // don't measure if any buttons are open which could distort measurement
    while (holdButton->isOpen() || rainGauge->isOpen())
    {
        return;
    }
    tempSensor->measure();
    tempSensor->sendTLVPacket();
    digitalWrite(LED_BLUE, 0);
}

/* drop-in replacement for `setup()` from arduino core */
void customSetup()
{
    delay(1000); // wait for serial connection to pick up first
    tempSensor->measure(); // doesn't matter if the reading is incorrect, we just need a reference
                           // point for memory allocation
    serialTLV->sendHardReset();
}

/* drop-in replacement for `loop()` in arduino code */
void customLoop()
{
    if (!paused)
    {
        if (rainGauge->isPressed())
            handleRainGauge();
        if (tempTimer->ready())
            handleMeasureTemp();
    }
    if (holdButton->isPressed())
        handlePause();
    digitalWrite(LED_GREEN, (millis() / 5000) % 2);  // blink the LED every 5 seconds
}

/* vendored setup for samd chip */
void arduinoCoreSamdMain()
{
    init();
    __libc_init_array();
    initVariant();
    delay(1);
#if defined(USBCON)
    USBDevice.init();
    USBDevice.attach();
#endif
}

int main()
{
    arduinoCoreSamdMain();
    customSetup();

    /* the main loop */
    for (;;)
    {
        customLoop();
        if (serialEventRun)
            serialEventRun();
    }
    return 0;
}
