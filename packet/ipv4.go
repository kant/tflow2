package packet

import (
	"fmt"
	"unsafe"
)

var (
	SizeOfIPv4Header = unsafe.Sizeof(IPv4Header{})
)

type IPv4Header struct {
	DstAddr             [4]byte
	SrcAddr             [4]byte
	HeaderChecksum      uint16
	Protocol            uint8
	TTL                 uint8
	FlagsFragmentOffset uint16
	Identification      uint16
	TotalLength         uint16
	DSCP                uint8
	VersionHeaderLength uint8
}

func DecodeIPv4(raw unsafe.Pointer, length uint32) (*IPv4Header, error) {
	if SizeOfIPv4Header > uintptr(length) {
		return nil, fmt.Errorf("Frame is too short: %d", length)
	}

	return (*IPv4Header)(unsafe.Pointer(uintptr(raw) - SizeOfIPv4Header)), nil
}
