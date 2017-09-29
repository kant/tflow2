package packet

import (
	"fmt"
	"unsafe"
)

var (
	SizeOfIPv6Header = unsafe.Sizeof(IPv6Header{})
)

type IPv6Header struct {
	DstAddr                      [16]byte
	SrcAddr                      [16]byte
	HotLimit                     uint8
	NextHeader                   uint8
	PayloadLength                uint16
	VersionTrafficClassFlowLabel uint32
}

func DecodeIPv6(raw unsafe.Pointer, length uint32) (*IPv6Header, error) {
	if SizeOfIPv6Header > uintptr(length) {
		return nil, fmt.Errorf("Frame is too short: %d", length)
	}

	return (*IPv6Header)(unsafe.Pointer(uintptr(raw) - SizeOfIPv6Header)), nil
}
