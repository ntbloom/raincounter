/* all-or-nothing, totally unforgiving unit-test suite for tlv encryption */

#include <cassert>
#include <iostream>

#include "../src/tlv.hpp"

#define OK cout << __func__ << "...OK\n"

using namespace std;
using tlv::TLV;

void assert_equal_arrays(int len, unsigned char one[], unsigned char two[]) {
    for (int i = 0; i < len - 1; i++) {
        assert(one[i] == two[i]);
    }
}

/* dummy test */
void test_variable_length() {
    assert(sizeof(float) == 4);  // value on arduino
    assert(sizeof(int) == 4);    // value on arduino
    OK;
}

/* rain event */
void test_tlv_rain_packet() {
    unsigned char t, l, v;
    t = 0;
    l = 1;
    v = 1;
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, v};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(3, expected, actual);
    delete tlv;
    OK;
}

/* soft-reset event */
void test_tlv_soft_reset_packet() {
    unsigned char t, l, v;
    t = 2;
    l = 1;
    v = 1;
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, v};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(3, expected, actual);
    delete tlv;
    OK;
}

/* hard-reset event */
void test_tlv_hard_reset_packet() {
    unsigned char t, l, v;
    t = 3;
    l = 1;
    v = 1;
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, v};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(3, expected, actual);
    delete tlv;
    OK;
}

/* happy path for temperature measurement */
void test_tlv_positive_temperature_packet() {
    unsigned char t, l;
    int v;
    t = 1;
    l = 4;
    v = 24;  // 24C, or 75F
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, 0, 0, 1, 8};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(6, expected, actual);
    delete tlv;
    OK;
}

/* zero temperature value */
void test_tlv_zero_temperature_packet() {
    unsigned char t, l;
    int v;
    t = 1;
    l = 4;
    v = 0;  // 0C or 32F
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, 0, 0, 0, 0};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(6, expected, actual);
    delete tlv;
    OK;
}

/* negative temperature values */
void test_tlv_negative_temperature_packet() {
    unsigned char t, l;
    int v;
    t = 1;
    l = 4;
    v = -24;  // -24C, or -11F
    TLV* tlv = new TLV(t, v);
    unsigned char expected[] = {t, l, 15, 15, 14, 7};
    unsigned char* actual = tlv->encode();
    assert_equal_arrays(6, expected, actual);
    delete tlv;
    OK;
}

int main() {
    test_variable_length();
    test_tlv_rain_packet();
    test_tlv_soft_reset_packet();
    test_tlv_hard_reset_packet();
    test_tlv_positive_temperature_packet();
    test_tlv_zero_temperature_packet();
    test_tlv_negative_temperature_packet();
}
