#ifndef _LED_MANAGER_HPP_
#define _LED_MANAGER_HPP_

namespace components
{

/* is the LED turned on or off */
constexpr int LED_ON = 1;
constexpr int LED_OFF = 0;

class Light
{
  private:
    unsigned long _duration;
    int _pin;
    unsigned long _start;
    bool _permanentlyOn;
    int _state;

  public:
    /* @param duration: how long to be on in milliseconds */
    Light(int pin);

    /* turn on the LED  for duration */
    void flash(void);

    /* turn on the LED indefinitely */
    void on(void);

    /* turn off the LED indefinitely */
    void off(void);

    /* check if it's time to turn off the light */
    void check(void);
};
}; // namespace components

#endif
