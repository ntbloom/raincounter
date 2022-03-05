#include "temp36.hpp"
#include <Arduino.h>

using namespace components;
using namespace tlv;

Temp36::Temp36(int pin, float voltage)
{
    _pin = pin;
    pinMode(_pin, INPUT);
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
    analogReadResolution(16);
    int reading = analogRead(_pin);
    float intermed = reading * _voltage / 0xffff;
    float offset = 0.5;
    float tempC = (intermed - offset) * 100; // 10mv per degree with 500 mV offset

    _valC = (int)tempC;
}

/* send message over serial port */
void Temp36::sendTLVPacket()
{
    _serialTLV->update(_valC);
    _serialTLV->sendPacket();
}
