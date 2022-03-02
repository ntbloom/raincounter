#include "component_serial_tlv.hpp"

#include "iserial_tlv.hpp"

using namespace tlv;

ComponentSerialTLV::ComponentSerialTLV(unsigned char tag, int value)
{
    // set the value as the initial value; we really only care about the size for TLV
    _tlv = new TLV(tag, value);
}

/* update the value of the tlv packet */
void ComponentSerialTLV::update(int value)
{
    _tlv->updateValue(value);
}

/* send a TLV packet encoded in SERIALFMT*/
void ComponentSerialTLV::sendPacket()
{
    _send(_tlv->encode(), SERIALFMT);
}
