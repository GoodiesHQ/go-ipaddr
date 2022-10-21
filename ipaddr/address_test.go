package ipaddr

import (
	"testing"
)

// Create an IP4Address object from a slice of bytes, does not raise an error
func TestCreateIP4AddressFromBytes(t *testing.T) {
	value := []byte{192, 168, 1, 123}
	addr, err := CreateIP4AddressFromBytes(value)

	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create proper IP4Address")
	}
}

// create an IP address from a Uint32, does not raise an error
func TestCreateIP4AddressFromUint32(t *testing.T) {
	value := uint32(0xc0_a8_01_7b)
	addr, err := CreateIP4AddressFromUint32(value)

	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create proper IP4Address")
	}
}

func TestCreateIP4Address(t *testing.T) {
	var addr IP4Address
	var err error

	// test IP only
	addr, err = CreateIP4Address("192.168.1.123")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create proper IP4Address")
	}

	// test IP with CIDR
	addr, err = CreateIP4Address("192.168.1.123/24")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create proper IP4Address")
	}

	// test IP with netmask
	addr, err = CreateIP4Address("192.168.1.123/255.255.255.0")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create proper IP4Address")
	}

	// test an invalid IP address
	addr, err = CreateIP4Address("100.200.300.400")
	if err == nil {
		t.Errorf("failed to return an error on invalid IP4Address")
	}

}

// Determines if the IP address is a valid host mask
func TestIsValidNetMask(t *testing.T) {
	var addr IP4Address
	var err error

	addr, err = CreateIP4Address("255.255.224.0")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if !addr.IsValidNetMask() {
		t.Errorf("failed to correctly identify a valid netmask")
	}

	addr, err = CreateIP4Address("255.255.0.224")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.IsValidNetMask() {
		t.Errorf("failed to correctly identify an invalid netmask")
	}
}

// Returns a netmask IP address into a CIDR. Returns -1 if it is not a valid netmask
func TestToCIDR(t *testing.T) {
	var addr IP4Address
	var err error

	addr, err = CreateIP4Address("0.0.0.0")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.ToCIDR() != 0 {
		t.Errorf("failed to correctly create the CIDR")
	}

	addr, err = CreateIP4Address("255.255.255.255")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.ToCIDR() != 32 {
		t.Errorf("failed to correctly create the CIDR")
	}

	addr, err = CreateIP4Address("255.255.224.0")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.ToCIDR() != 19 {
		t.Errorf("failed to correctly create the CIDR")
	}

	addr, err = CreateIP4Address("255.255.0.224")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.ToCIDR() != -1 {
		t.Errorf("failed to correctly create the CIDR")
	}

}

func TestToUint32(t *testing.T) {
	var addr IP4Address
	var err error

	addr, err = CreateIP4Address("192.168.1.123")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if addr.ToUint32() != 0xc0_a8_01_7b {
		t.Errorf("failed to create the correct uint32 value")
	}
}

func TestAddressString(t *testing.T) {
	var addr IP4Address
	var err error

	addr, err = CreateIP4Address("192.168.1.123")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}
	if addr.String() != "192.168.1.123" {
		t.Errorf("failed to create the correct string value")
	}
}
