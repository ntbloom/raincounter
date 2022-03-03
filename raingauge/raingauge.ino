/* arduino_rain_gauge.ino
 *
 * Increment rain gauge counter on click, display on screen
 *
 */

#include <Arduino.h>

/* specify the board to be used */
#include "src/boards/nano33ble.cpp"

void ledSetup(void)
{
    pinMode(LED_GREEN, OUTPUT);
    pinMode(LED_RED, OUTPUT);
    pinMode(LED_BLUE, OUTPUT);
}

void setup()
{
    ledSetup();
}

void loop()
{
    // digitalWrite(LED_RED, (millis() / 1000) % 2);   // blink the LED every 5 seconds
    // digitalWrite(LED_BLUE, (millis() / 1000) % 2);  // blink the LED every 5 seconds
    digitalWrite(LED_GREEN, (millis() / 1000) % 2); // blink the LED every 5 seconds
}

