#include "static_serial_tlv.hpp"

#include "iserial_tlv.hpp"

#define RAIN_COUNTER 0
#define SOFT_RESET 2
#define HARD_RESET 3
#define PAUSE 4
#define UNPAUSE 5

using namespace tlv;
StaticSerialTLV::StaticSerialTLV()
{
    _rain = _makeTLV(RAIN_COUNTER);
    _softReset = _makeTLV(SOFT_RESET);
    _hardReset = _makeTLV(HARD_RESET);
    _pause = _makeTLV(PAUSE);
    _unpause = _makeTLV(UNPAUSE);
}

/* preallocate a TLV packet */
unsigned char *StaticSerialTLV::_makeTLV(unsigned char tag)
{
    unsigned char val = 1;
    TLV *tlv = new TLV(tag, val);
    return tlv->encode();
}

/* indicates a rain gauge tipper was incremented */
void StaticSerialTLV::sendRainEvent()
{
    _send(_rain, SERIALFMT);
}

/* indicate soft reset (rain counter reset) just happened */
void StaticSerialTLV::sendSoftReset()
{
    _send(_softReset, SERIALFMT);
}

/* send right after boot, indicating a hard reset happened */
void StaticSerialTLV::sendHardReset()
{
    _send(_hardReset, SERIALFMT);
}

/* send when sensor is paused */
void StaticSerialTLV::sendPause()
{
    _send(_pause, SERIALFMT);
}

/* send when sensor is unpaused */
void StaticSerialTLV::sendUnpause()
{
    _send(_unpause, SERIALFMT);
}

StaticSerialTLV::~StaticSerialTLV()
{
    delete[] _softReset;
    delete[] _hardReset;
    delete[] _pause;
    delete[] _unpause;
}
