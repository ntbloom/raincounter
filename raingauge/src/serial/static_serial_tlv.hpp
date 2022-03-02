#ifndef _STATIC_SERIAL_TLV_HPP_
#define _STATIC_SERIAL_TLV_HPP_

#include "iserial_tlv.hpp"

namespace tlv
{
/* Manage static TLV packets that get allocated once and reused */
class StaticSerialTLV : public ISerialTLV
{
  private:
    unsigned char *_rain, *_softReset, *_hardReset, *_pause, *_unpause;
    unsigned char *_makeTLV(unsigned char tag);

  public:
    StaticSerialTLV();
    void sendHex();
    void sendRainEvent();
    void sendSoftReset();
    void sendHardReset();
    void sendPause();
    void sendUnpause();
    ~StaticSerialTLV();
};
}; // namespace tlv
#endif
