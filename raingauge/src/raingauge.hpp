#ifndef _RAINGAUGE_HPP_
#define _RAINGAUGE_HPP_

#include "button.hpp"

/* tipper rain gauge */
namespace components
{
class Raingauge : public Button
{
  private:
    float _mmPerCount;
    unsigned long _count;
    float _mmTotal;
    void _updateValues();

  public:
    Raingauge(int pin, unsigned long msDelay, float mmPerCount);
    void addCount();
    void resetCount();
};
}; // namespace components

#endif
