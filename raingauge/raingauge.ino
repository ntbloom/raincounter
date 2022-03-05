/* raingauge.ino
 *
 *  Measure the rain and temperature and send data over serial port as TLV
 *
 */

#include "src/boards/nano33ble.hpp"
#include "src/config.hpp"
#include "src/led_manager.hpp"
#include "src/raingauge.hpp"
#include "src/serial/static_serial_tlv.hpp"
#include "src/temp36.hpp"
#include "src/timer.hpp"
#include <Arduino.h>

/* Measure Temperature */
static components::Temp36 tempSensor(TEMP_PIN, TEMP_VOLTAGE);
static utilities::Timer tempTimer(TEMP_INTERVAL_SEC);

/* Raingauge */
static components::Raingauge raingauge(RAIN_PIN, RAINGAUGE_DEBOUNCE_MS, RAINGAUGE_AMT_MM);

/* Serial line */
tlv::StaticSerialTLV *serialTLV = new tlv::StaticSerialTLV();

/* Pause button */
static components::Button pauseButton(PAUSE_PIN, BUTTON_DEBOUNCE_MS, HIGH);
static bool PAUSED = false;

/* LED Behavior */
static components::Light greenLED(LED_TEMP_GREEN);
static components::Light redLED(LED_PAUSE_RED);
static components::Light blueLED(LED_RAIN_BLUE);

/* check the LEDs to see if we need to turn any off */
void checkLED(void)
{
    greenLED.check();
    redLED.check();
    blueLED.check();
}

/* measure temperature and send data to serial port */
void handleMeasureTemp(void)
{
    // if a switch is open we could get a distorted reading, skip it for now
    if (pauseButton.isOpen() || raingauge.isOpen())
    {
        return;
    }
    tempSensor.measure();
    tempSensor.sendTLVPacket();
    greenLED.flash();
}

/* count a rain event */
void handleRainEvent(void)
{
    raingauge.addCount();
    serialTLV->sendRainEvent();
    blueLED.flash();
}

/* stop all measurements and turn on red pause LED */
void handlePause(void)
{
    if (PAUSED)
    {
        PAUSED = false;
        serialTLV->sendUnpause();
        redLED.off();
    }
    else
    {
        PAUSED = true;
        serialTLV->sendPause();
        redLED.on();
    }
}

void setup()
{
    // give a chance for the serial port to pick up
    delay(1000);
    serialTLV->sendHardReset();
}

void loop()
{
    checkLED();
    if (!PAUSED)
    {
        if (raingauge.isPressed())
        {
            handleRainEvent();
        }

        if (tempTimer.ready())
        {
            handleMeasureTemp();
        }
    }

    if (pauseButton.isPressed())
    {
        handlePause();
    }
}
