#include "led_manager.hpp"
#include "Arduino.h"
#include "config.hpp"

using namespace components;

Light::Light(int pin)

{
    _pin = pin;
    pinMode(_pin, OUTPUT);
    _start = millis();
    _permanentlyOn = false;
}

void Light::flash(void)
{
    _state = LED_ON;
    digitalWrite(_pin, LED_ON);
    _start = millis();
}

void Light::on(void)
{
    _state = LED_ON;
    digitalWrite(_pin, LED_ON);
    _permanentlyOn = true;
}

void Light::off(void)
{
    _state = LED_OFF;
    digitalWrite(_pin, LED_OFF);
    _permanentlyOn = false;
}

void Light::check(void)
{
    if (_state == LED_OFF || _permanentlyOn)
    {
        return;
    }

    if ((millis() - _start) >= LED_FLASH_DURATION_MS || (millis() < _start))
    {
        digitalWrite(_pin, LED_OFF);
        _state = LED_OFF;
    }
}
