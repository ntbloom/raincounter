#ifndef _TLV_HPP_
#define _TLV_HPP_

namespace tlv
{
class TLV
{
  private:
    unsigned char *_payload;

  public:
    TLV(unsigned char tag, unsigned char value);
    TLV(unsigned char tag, int value);
    unsigned char *encode();
    void updateValue(int value);
    ~TLV();
};
} // namespace tlv

#endif
