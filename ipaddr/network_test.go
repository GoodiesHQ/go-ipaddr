package ipaddr

import (
	"testing"
)

func TestIterate(t *testing.T) {
	var net IP4Network
	var netIt *IP4Network
	var nets []IP4Network
	var it IP4NetworkIterator
	var err error

	net, err = CreateIP4Network("192.168.1.123/29")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	nets = make([]IP4Network, 0)
	it = net.Iterate(true)
	for netIt = it.Next(); netIt != nil; netIt = it.Next() {
		nets = append(nets, *netIt)
	}

	if len(nets) != 6 {
		t.Errorf("failed to create correct number of IP4Network objects")
	}

	if nets[0].Address.String() != "192.168.1.121" ||
		nets[1].Address.String() != "192.168.1.122" ||
		nets[2].Address.String() != "192.168.1.123" ||
		nets[3].Address.String() != "192.168.1.124" ||
		nets[4].Address.String() != "192.168.1.125" ||
		nets[5].Address.String() != "192.168.1.126" {

		t.Errorf("failed to create correct array of IP4Network objects: %v", nets[0].Address.String())
	}

	nets = make([]IP4Network, 0)
	it = net.Iterate(false)
	for netIt = it.Next(); netIt != nil; netIt = it.Next() {
		nets = append(nets, *netIt)
	}

	if len(nets) != 8 {
		t.Errorf("failed to create correct number of IP4Network objects")
	}

	if nets[0].Address.String() != "192.168.1.120" ||
		nets[1].Address.String() != "192.168.1.121" ||
		nets[2].Address.String() != "192.168.1.122" ||
		nets[3].Address.String() != "192.168.1.123" ||
		nets[4].Address.String() != "192.168.1.124" ||
		nets[5].Address.String() != "192.168.1.125" ||
		nets[6].Address.String() != "192.168.1.126" ||
		nets[7].Address.String() != "192.168.1.127" {

		t.Errorf("failed to create correct array of IP4Network objects")
	}

}

func TestNetworkString(t *testing.T) {
	var net IP4Network
	var err error

	net, err = CreateIP4Network("192.168.1.123/255.255.224.0")
	if err != nil {
		t.Errorf("failed to create network")
	}

	if net.String() != "192.168.1.123/19" {
		t.Errorf("invalid network string")
	}
}

func TestNetworkContains(t *testing.T) {
	var net IP4Network
	var addr IP4Address
	var err error

	net, err = CreateIP4Network("192.168.1.123/19")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	addr, err = CreateIP4Address("192.168.2.34")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if !net.Contains(addr) {
		t.Errorf("failed to correctly determine if network contains address")
	}

	addr, err = CreateIP4Address("192.168.100.123")
	if err != nil {
		t.Errorf("failed to create IP4Address")
	}

	if net.Contains(addr) {
		t.Errorf("failed to correctly determine if network contains address")
	}
}

func TestNetworkAddress(t *testing.T) {
	var net IP4Network
	var addr IP4Address
	var err error

	net, err = CreateIP4Network("192.168.123.123/19")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	addr = net.NetworkAddress()
	if addr.String() != "192.168.96.0" {
		t.Errorf("failed to correctly determine the network address")
	}

	net, err = CreateIP4Network("192.168.12.34/255.255.248.0")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	addr = net.NetworkAddress()
	if addr.String() != "192.168.8.0" {
		t.Errorf("failed to correctly determine the network address")
	}
}

func TestBroadcastAddress(t *testing.T) {
	var net IP4Network
	var addr IP4Address
	var err error

	net, err = CreateIP4Network("192.168.123.123/19")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	addr = net.BroadcastAddress()
	if addr.String() != "192.168.127.255" {
		t.Errorf("failed to correctly determine the network address")
	}

	net, err = CreateIP4Network("192.168.12.34/255.255.248.0")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	addr = net.BroadcastAddress()
	if addr.String() != "192.168.15.255" {
		t.Errorf("failed to correctly determine the network address")
	}
}

func TestCreateIP4Netmask(t *testing.T) {
	var addr IP4Address
	var err error

	addr, err = CreateIP4Netmask("19")
	if err != nil {
		t.Errorf("failed to create IP4Netmask")
	}

	if addr.ToCIDR() != 19 || addr.String() != "255.255.224.0" {
		t.Errorf("failed to create valid IP4Netmask")
	}

	addr, err = CreateIP4Netmask("255.255.248.0")
	if err != nil {
		t.Errorf("failed to create IP4Netmask")
	}

	if addr.ToCIDR() != 21 || addr.String() != "255.255.248.0" {
		t.Errorf("failed to create valid IP4Netmask")
	}

	_, err = CreateIP4Netmask("1.2.3.4")
	if err == nil {
		t.Errorf("failed to raise error on invalid IP4Netmask")
	}

	_, err = CreateIP4Netmask("35")
	if err == nil {
		t.Errorf("failed to raise error on invalid IP4Netmask")
	}
}

func TestCreateIP4Network(t *testing.T) {
	var net IP4Network
	var err error

	net, err = CreateIP4Network("192.168.1.123/21")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	if net.Address.String() != "192.168.1.123" || net.Netmask.String() != "255.255.248.0" {
		t.Errorf("failed to create valid IP4Network")
	}

	net, err = CreateIP4Network("192.168.1.123/255.255.248.0")
	if err != nil {
		t.Errorf("failed to create IP4Network")
	}

	if net.Address.String() != "192.168.1.123" || net.Netmask.String() != "255.255.248.0" {
		t.Errorf("failed to create valid IP4Network")
	}

	_, err = CreateIP4Network("100.200.300.400/255.255.248.0")
	if err == nil {
		t.Errorf("failed to raise error on invalid IP4Network")
	}

	_, err = CreateIP4Network("192.168.1.123/255.255.0.248")
	if err == nil {
		t.Errorf("failed to raise error on invalid IP4Network")
	}

	_, err = CreateIP4Network("192.168.1.123/35")
	if err == nil {
		t.Errorf("failed to raise error on invalid IP4Network")
	}
}
