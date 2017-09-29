package packet

import (
	"fmt"
	"unsafe"
)

const (
	UDP = 17
)

var (
	SizeOfUDPHeader = unsafe.Sizeof(UDPHeader{})
)

type UDPHeader struct {
	Checksum uint16
	Length   uint16
	DstPort  uint16
	SrcPort  uint16
}

func DecodeUDP(raw unsafe.Pointer, length uint32) (*UDPHeader, error) {
	if SizeOfTCPHeader > uintptr(length) {
		return nil, fmt.Errorf("Frame is too short: %d", length)
	}

	return (*UDPHeader)(unsafe.Pointer(uintptr(raw) - SizeOfUDPHeader)), nil
}
