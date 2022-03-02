#ifndef _LED_MANAGER_HPP_
#define _LED_MANAGER_HPP_

namespace components
{

typedef int led_t;
typedef int led_duration_t;

/* is the LED turned on or off */
enum class LED_State
{
    LED_ON = 1,
    LED_OFF = 0
};

/* control an individual LED */
struct Led
{
  public:
    /* initialize an LED at a GPIO pin to turn on for a given duration */
    Led(led_t pin, led_duration_t duration);
    void turnOn();

  private:
    led_t _pin;
    int _start;
    led_duration_t _duration;
    LED_State _state;
};

struct LedManager
{
  private:
    /* set pins for status, measurement, and pause LEDS */
    Led _status, _measure, _pause;

  public:
    LedManager(led_t status, led_t measure, led_t pause);

    /* trigger the pin on */
    void trigger(led_t pin);
};
} // namespace components

#endif
