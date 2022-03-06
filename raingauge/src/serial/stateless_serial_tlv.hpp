#ifndef _STATEFUL_SERIAL_TLV_HPP_
#define _STATEFUL_SERIAL_TLV_HPP_

#include "iserial_tlv.hpp"

namespace tlv
{
/* Manage static TLV packets that get allocated once and reused */
class StatelessSerialTLV : public ISerialTLV
{
  private:
    unsigned char *_rain, *_softReset, *_hardReset, *_pause, *_unpause;
    unsigned char *_makeTLV(unsigned char tag);

  public:
    StatelessSerialTLV();
    void sendHex();
    void sendRainEvent();
    void sendSoftReset();
    void sendHardReset();
    void sendPause();
    void sendUnpause();
    ~StatelessSerialTLV();
};
}; // namespace tlv
#endif
