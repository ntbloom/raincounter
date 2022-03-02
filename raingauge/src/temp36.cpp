#include "temp36.hpp"

using namespace components;
using namespace tlv;

Temp36::Temp36(int pin, float voltage)
{
    _pin = pin;
    _voltage = voltage;
    measure();

    // allocate memory for the TLV packet
    _serialTLV = new ComponentSerialTLV(_tag, _valC);
}

/* calculate the temperature, store values in memory */
void Temp36::measure()
{
    /* note, values could get distorted based on voltage flow around the board
     * put logic into scripts for __when__ to measure to make sure values are accurate
     */
    analogReadResolution(12);
    int reading = analogRead(_pin);
    float intermed = reading * _voltage / 4096;
    float tempC = (intermed - 0.5) * 100; // 10mv per degree with 500 mV offset
    float tempF = (tempC * 9.0 / 5.0) + 32.0;

    _valF = (int)tempF;
    _valC = (int)tempC;
    _tempF = String(_valF) + "F";
    _tempC = String(_valC) + "C";
}

/* getter for F, as string */
String Temp36::tempF()
{
    return _tempF;
}

/* getter for C, as string */
String Temp36::tempC()
{
    return _tempC;
}

/* send message over serial port */
void Temp36::sendTLVPacket()
{
    _serialTLV->update(_valC);
    _serialTLV->sendPacket();
}
