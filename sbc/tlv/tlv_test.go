package tlv_test

import (
	"fmt"
	"testing"

	"github.com/ntbloom/rainbase/pkg/tlv"
)

// Test static packets like rain event, etc.
func TestStaticVals(t *testing.T) {
	exp := 1
	rain := []byte{48, 49, 49}
	if !verifyValToInt(rain, exp) {
		t.Fail()
	}
	softReset := []byte{50, 49, 49}
	if !verifyValToInt(softReset, exp) {
		t.Fail()
	}
	hardReset := []byte{51, 49, 49}
	if !verifyValToInt(hardReset, exp) {
		t.Fail()
	}
	pause := []byte{52, 49, 49}
	if !verifyValToInt(pause, exp) {
		t.Fail()
	}
}

// Test temperature packets
func Test18(t *testing.T) {
	// positives
	temp1 := []byte{49, 52, 48, 48, 49, 50, 10}
	exp := 18
	if !verifyValToInt(temp1, exp) {
		t.Fail()
	}
}
func Test25(t *testing.T) {
	temp2 := []byte{49, 52, 48, 48, 49, 57, 10}
	exp := 25
	if !verifyValToInt(temp2, exp) {
		t.Fail()
	}
}
func Test26(t *testing.T) {
	temp3 := []byte{49, 52, 48, 48, 49, 65, 10}
	exp := 26
	if !verifyValToInt(temp3, exp) {
		t.Fail()
	}
}

func Test0(t *testing.T) {
	// zero
	temp4 := []byte{49, 52, 48, 48, 48, 48, 10}
	exp := 0
	if !verifyValToInt(temp4, exp) {
		t.Fail()
	}
}

func TestMinus24(t *testing.T) {
	// negatives
	temp5 := []byte{49, 52, 70, 70, 69, 55, 10}
	exp := -24
	if !verifyValToInt(temp5, exp) {
		t.Fail()
	}
}

func verifyValToInt(raw []byte, expected int) bool {
	tlvPacket, err := tlv.NewTLV(raw)
	if tlvPacket == nil {
		fmt.Printf("error making tlvPacket: %s\n", err)
		return false
	}
	actual := tlvPacket.Value
	if actual != expected {
		fmt.Printf("expected=%d, actual=%d\n", expected, actual)
		return false
	}
	return true
}
