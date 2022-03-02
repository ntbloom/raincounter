#ifndef _BUTTON_HPP_
#define _BUTTON_HPP_

#include "Arduino.h"

/* basic class for controlling and debouncing buttons */
namespace components
{
class Button
{
  private:
    int _pin;
    unsigned long _msDelay;
    unsigned long _lastDebounce = 0;
    unsigned long _now = 0;
    bool _high;
    bool _read = LOW;

  public:
    Button(int pin, unsigned long msDelay, bool high);
    bool isOpen();
    bool isPressed();
};
}; // namespace components

#endif
