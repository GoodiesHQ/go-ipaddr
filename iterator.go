package ipaddr

import (
	"fmt"
)

type IP4NetworkIteration interface {
	Next() IP4Network
	Inc() error
	Dec() error
}

type IP4NetworkIterator struct {
	HostsOnly        bool   // determines if Network and Broadcast address should be given
	started          bool   // if the iterator has been started
	address          uint32 // IP4Address address object converted to uint32
	netmask          uint32 // IP4Address netmask object to iterate over (changes each iteration)
	networkAddress   uint32 // IP4Address network address (calculated first IP)
	broadcastAddress uint32 // IP4Address broadcast address (calculated last IP)
}

func (iter *IP4NetworkIterator) Inc() error {
	if iter.address == iter.broadcastAddress {
		return fmt.Errorf("no more addresses")
	}
	iter.address += 1
	return nil
}

func (iter *IP4NetworkIterator) Dec() error {
	if iter.address == iter.networkAddress {
		return fmt.Errorf("no more addresses")
	}
	iter.address -= 1
	return nil
}

func (iter *IP4NetworkIterator) Next() *IP4Network {
	if !iter.started {
		iter.started = true
		if iter.address == iter.networkAddress && iter.HostsOnly {
			iter.Inc()
		}
	} else if iter.started && iter.address != iter.broadcastAddress {
		iter.Inc()
	} else if iter.address == iter.broadcastAddress {
		return nil
	}

	if iter.address == iter.broadcastAddress {
		if iter.HostsOnly {
			return nil
		}
	}

	address, _ := CreateIP4AddressFromUint32(iter.address)
	netmask, _ := CreateIP4AddressFromUint32(iter.netmask)

	return &IP4Network{
		Address: address,
		Netmask: netmask,
	}
}
