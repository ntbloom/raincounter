#include "led_manager.hpp"

using namespace components;

LED_State &operator++(LED_State &t)
{
    switch (t)
    {
        case LED_State::LED_ON:
            return t = LED_State::LED_OFF;
        case LED_State::LED_OFF:
            return t = LED_State::LED_ON;
    }
}

Led::Led(led_t pin, led_duration_t duration)
{
    // do something!
}
