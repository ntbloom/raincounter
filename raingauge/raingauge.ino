/* raingauge.ino
 *
 *  Measure the rain and temperature and send data over serial port as TLV
 *
 */

#include "src/boards/nano33ble.hpp"
#include "src/config.hpp"
#include "src/led_manager.hpp"
#include "src/raingauge.hpp"
#include "src/serial/stateless_serial_tlv.hpp"
#include "src/temp36.hpp"
#include "src/timer.hpp"
#include <Arduino.h>

/* Measure Temperature */
static components::Temp36 tempSensor(TEMP_PIN, TEMP_VOLTAGE);
static utilities::Timer tempTimer(TEMP_INTERVAL_SEC);

/* Raingauge */
static components::Raingauge raingauge(RAIN_PIN, RAINGAUGE_DEBOUNCE_MS, RAINGAUGE_AMT_MM);

/* Serial line */
tlv::StatelessSerialTLV *serialTLV = new tlv::StatelessSerialTLV();

/* Pause button */
static components::Button pauseButton(PAUSE_PIN, BUTTON_DEBOUNCE_MS, HIGH);
static bool PAUSED = false;

/* LED Behavior */
static components::Light greenLED(LED_TEMP_GREEN);
static components::Light redLED(LED_PAUSE_RED);
static components::Light blueLED(LED_RAIN_BLUE);
static components::Light lights[]{greenLED, redLED, blueLED};

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
    // wait a bit for serial port to pick up, blink lights in the meantime
    for (int i = 0; i < 3; i++)
    {
        for (auto &light : lights)
        {
            light.on();
            delay(100);
            light.off();
        }
    }
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
