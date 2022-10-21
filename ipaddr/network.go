package ipaddr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const IP4CIDRRegex = `^([0-9]|1[0-9]|2[0-9]|3[0-2])$`
const BitsPerByte = 8

type IP4Network struct {
	Address IP4Address
	Netmask IP4Address
}

func (net *IP4Network) Iterate(hostsOnly bool) IP4NetworkIterator {
	networkAddress := net.NetworkAddress()
	broadcastAddress := net.BroadcastAddress()

	return IP4NetworkIterator{
		HostsOnly:        hostsOnly,
		started:          false,
		address:          networkAddress.ToUint32(),
		netmask:          net.Netmask.ToUint32(),
		networkAddress:   networkAddress.ToUint32(),
		broadcastAddress: broadcastAddress.ToUint32(),
	}
}

func (net IP4Network) String() string {
	return fmt.Sprintf("%s/%d", net.Address.String(), net.Netmask.ToCIDR())
}

func (net *IP4Network) Contains(addr IP4Address) bool {
	netNetworkAddress := net.NetworkAddress()
	netBroadcastAddress := net.BroadcastAddress()

	netNetworkValue := netNetworkAddress.ToUint32()
	netBroadcastValue := netBroadcastAddress.ToUint32()
	addrValue := addr.ToUint32()

	return (addrValue >= netNetworkValue) && (addrValue <= netBroadcastValue)
}

func (net *IP4Network) NetworkAddress() IP4Address {
	var addr IP4Address

	// calculate the lowest IP value (network address)
	for i := range net.Address {
		addr[i] = net.Address[i] & net.Netmask[i]
	}
	return addr
}

func (net *IP4Network) NetworkNetwork() IP4Network {
	return IP4Network{
		Address: net.NetworkAddress(),
		Netmask: net.Netmask,
	}
}

func (net *IP4Network) BroadcastAddress() IP4Address {
	var addr IP4Address

	// calculate the highest IP value (broadcast address)
	for i := range net.Address {
		addr[i] = net.Address[i] | (^net.Netmask[i])
	}

	return addr
}

func (net *IP4Network) BroadcastNetwork() IP4Network {
	return IP4Network{
		Address: net.BroadcastAddress(),
		Netmask: net.Netmask,
	}
}

func CreateIP4Netmask(netmask string) (IP4Address, error) {
	netmask = strings.TrimSpace(netmask)

	// Check if a CIDR is provided, 0-32
	if match, _ := regexp.MatchString(IP4CIDRRegex, netmask); match {
		// convert the CIDR string to an integer
		cidr, err := strconv.Atoi(netmask)
		if err != nil {
			return IP4Address{}, err
		}

		// 32 bits all 1's
		var value uint32 = 0xffffffff

		// Shift the value by the number of host bits (32 - network bits)
		value <<= ((IP4AddressSize * BitsPerByte) - cidr)

		// create new address from the uint32
		address, err := CreateIP4AddressFromUint32(value)
		if err != nil {
			return IP4Address{}, err
		}

		return address, nil
	}

	// check if an IPv4 address is provided
	if match, _ := regexp.MatchString(IP4AddressRegex, netmask); match {
		// create an IP from the string
		address, err := CreateIP4Address(netmask)
		if err != nil {
			return IP4Address{}, err
		}

		// check if the address is an invalid netmask
		if !address.IsValidNetMask() {
			return IP4Address{}, fmt.Errorf("invalid netmask value")
		}

		return address, nil
	}

	// no match found
	return IP4Address{}, fmt.Errorf("invalid netmask value")
}

func CreateIP4Network(network string) (IP4Network, error) {
	network = strings.TrimSpace(network)

	// all network strings should have a slash to separate the IP from the netmask/CIDR
	if strings.Count(network, "/") != 1 {
		return IP4Network{}, fmt.Errorf("expected one CIDR or netmask for the network")
	}

	// split into the network and mask
	parts := strings.Split(network, "/")

	// ensure the IP address is valid
	address, err := CreateIP4Address(parts[0])
	if err != nil {
		return IP4Network{}, err
	}

	// ensure the netmask/CIDR is valid
	netmask, err := CreateIP4Netmask(parts[1])
	if err != nil {
		return IP4Network{}, err
	}

	return IP4Network{
		Address: address,
		Netmask: netmask,
	}, nil
}
