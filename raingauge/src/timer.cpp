#include "timer.hpp"

using namespace utilities;

/* how long to wait, in seconds */
Timer::Timer(unsigned long interval)
{
    _interval = interval;
}

/* has that amount of time passed? */
bool Timer::ready()
{
    _now = millis();
    if (_count > _now)
    { // millis() overflowed back to zero
        _count = 0;
    }
    if ((_now - _count) > (_interval * 1000))
    {
        _count = _now;
        return true;
    }
    return false;
}
