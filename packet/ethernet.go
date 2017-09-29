package packet

import (
	"fmt"
	"net"
	"unsafe"

	"github.com/taktv6/tflow2/convert"
)

const (
	// EtherTypeARP is Address Resolution Protocol EtherType value
	EtherTypeARP = 0x0806

	// EtherTypeIPv4 is Internet Protocol version 4 EtherType value
	EtherTypeIPv4 = 0x0800

	// EtherTypeIPv6 is Internet Protocol Version 6 EtherType value
	EtherTypeIPv6 = 0x86DD

	// EtherTypeLACP is Link Aggregation Control Protocol EtherType value
	EtherTypeLACP = 0x8809

	// EtherTypeIEEE8021Q is VLAN-tagged frame (IEEE 802.1Q) EtherType value
	EtherTypeIEEE8021Q = 0x8100
)

var (
	SizeOfEthernetII = unsafe.Sizeof(ethernetII{})
)

// EthernetHeader represents layer two IEEE 802.11
type EthernetHeader struct {
	SrcMAC    net.HardwareAddr
	DstMAC    net.HardwareAddr
	EtherType uint16
}

type ethernetII struct {
	EtherType uint16
	SrcMAC    [6]byte
	DstMAC    [6]byte
}

func DecodeEthernet(raw unsafe.Pointer, length uint32) (*EthernetHeader, error) {
	if SizeOfEthernetII > uintptr(length) {
		return nil, fmt.Errorf("Frame is too short: %d", length)
	}

	ptr := unsafe.Pointer(uintptr(raw) - SizeOfEthernetII)
	ethHeader := (*ethernetII)(ptr)

	srcMAC := ethHeader.SrcMAC[:]
	dstMAC := ethHeader.DstMAC[:]

	srcMAC = convert.Reverse(srcMAC)
	dstMAC = convert.Reverse(dstMAC)

	h := &EthernetHeader{
		SrcMAC:    net.HardwareAddr(srcMAC),
		DstMAC:    net.HardwareAddr(dstMAC),
		EtherType: ethHeader.EtherType,
	}

	return h, nil
}
