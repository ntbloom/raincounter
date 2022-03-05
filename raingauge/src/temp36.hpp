#ifndef _TEMP36_HPP_
#define _TEMP36_HPP_

#include "serial/component_serial_tlv.hpp"
#include "tlv.hpp"

using tlv::ComponentSerialTLV;

/* component for TMP36 analog temperature sensors */
namespace components
{
class Temp36
{
  private:
    int _pin;
    float _voltage;
    const unsigned char _tag = 1;
    int _valC;
    ComponentSerialTLV *_serialTLV;

  public:
    Temp36(int pin, float voltage);
    void measure();
    void sendTLVPacket();
};
}; // namespace components

#endif
