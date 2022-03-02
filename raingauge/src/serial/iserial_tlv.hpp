#ifndef _ISERIAL_TLV_HPP_
#define _ISERIAL_TLV_HPP_

#include "../tlv.hpp"
#include "Arduino.h"

#define BAUD 115200
#define SERIALFMT HEX

namespace tlv
{
/* Abstract base class for sending TLV packets over a serial port */
class ISerialTLV
{
  protected:
    const int _baud = BAUD;
    virtual void _send(unsigned char *packet, int base)
    {
        Serial.begin(_baud);
        for (unsigned char i = 0; i < packet[1] + 2; i++)
        {
            Serial.print(packet[i], base);
        }
        Serial.print('\n');
        Serial.end();
    }
};
}; // namespace tlv
#endif
