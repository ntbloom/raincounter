#include "raingauge.hpp"

using namespace components;

Raingauge::Raingauge(int pin, unsigned long msDelay, float mmPerCount) : Button(pin, msDelay, HIGH)
{
    _count = 0;
    _mmPerCount = mmPerCount;
    _updateValues();
}

/* add a click to the counter */
void Raingauge::addCount()
{
    _count++;
    _updateValues();
}

void Raingauge::_updateValues()
{
    _mmTotal = _count * _mmPerCount;
}

/* reset the counters */
void Raingauge::resetCount()
{
    _count = 0;
    _updateValues();
}

