package ipaddr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const IP4AddressSize = 4
const IP6AddressSize = 16

type IPVersion int

const (
	IPUK IPVersion = 0 // IP version unknown
	IPV4 IPVersion = 4 // IPv4
	IPV6 IPVersion = 6 // IPv6
)

type IP4Address [IP4AddressSize]byte
type IP6Address [IP6AddressSize]byte

type IPAddress interface {
	String() string     // returns a string representation of the ip address
	Version() IPVersion // returns the IP version of the ip address
}

func (addr *IP4Address) Version() IPVersion {
	return IPV4
}

func (addr *IP6Address) Version() IPVersion {
	return IPV6
}

const IP4AddressRegex = `^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`

// Create an IP4Address object from a slice of bytes
func CreateIP4AddressFromBytes(value []byte) (IP4Address, error) {
	if len(value) != IP4AddressSize {
		return IP4Address{}, fmt.Errorf("expected %d bytes, got %d", IP4AddressSize, len(value))
	}
	var address IP4Address
	copy(address[:], value[:])
	return address, nil
}

// create an IP address from a Uint32, does not raise an error as any uint32 can be converted to an IP
func CreateIP4AddressFromUint32(value uint32) (IP4Address, error) {
	bytes := make([]byte, IP4AddressSize)
	for i := 0; i < IP4AddressSize; i++ {
		bytes[IP4AddressSize-i-1] = byte(value & 0xff)
		value >>= 8
	}
	return CreateIP4AddressFromBytes(bytes)
}

// create an IP address from a string
func CreateIP4Address(address string) (IP4Address, error) {
	addr := IP4Address{}
	address = strings.TrimSpace(address)
	slashCount := strings.Count(address, "/")

	if slashCount > 1 {
		return IP4Address{}, fmt.Errorf("too many slashes in the provided address")
	}

	if slashCount > 0 {
		address = strings.Split(address, "/")[0]
	}

	matched, err := regexp.Match(IP4AddressRegex, []byte(address))
	if err != nil {
		return IP4Address{}, err
	}

	if !matched {
		return IP4Address{}, fmt.Errorf("invalid IP address provided")
	}

	if strings.Contains(address, "/") {
		address = strings.SplitN(address, "/", 1)[0]
	}

	for i, octet := range strings.Split(address, ".") {
		octetValue, _ := strconv.Atoi(octet)
		addr[i] = byte(octetValue)
	}

	return addr, nil
}

// Determines if the IP address is a valid host mask
func (addr *IP4Address) IsValidNetMask() bool {
	// flag to determine if it has hit a zero
	zeroes := false
	for _, b := range *addr {
		for mask := byte(0b_1000_0000); mask > 0; mask >>= 1 {
			if b&mask == 0 {
				// set the zeroes flag upon the first instance of a zero.
				zeroes = true
			} else {
				// if a 1 follows any number of zeroes, it is invalid
				if zeroes {
					return false
				}
			}
		}
	}
	// any unset bits follow all set bits
	return true
}

// Returns a netmask IP address into a CIDR. Returns -1 if it is not a valid netmask
func (addr *IP4Address) ToCIDR() int {
	// ensure the address is a valid netmask
	if !addr.IsValidNetMask() {
		return -1
	}

	var cidr int = 0
	// calculate the number of set network bits
	for _, b := range *addr {
		for mask := byte(0b_1000_0000); mask > 0; mask >>= 1 {
			if b&mask == 0 {
				return cidr
			}
			cidr += 1
		}
	}
	return cidr
}

func (addr *IP4Address) ToUint32() uint32 {
	var value uint32
	// rebuild uint32 from the slize of bytes
	for i := 0; i < IP4AddressSize; i++ {
		value <<= 8
		value |= uint32(addr[i])
	}
	return value
}

func (addr *IP4Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}
