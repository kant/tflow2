package packet

import (
	"fmt"
	"unsafe"
)

const (
	TCP = 6
)

var (
	SizeOfTCPHeader = unsafe.Sizeof(TCPHeader{})
)

type TCPHeader struct {
	UrgentPointer  uint16
	Checksum       uint16
	Window         uint16
	Flags          uint8
	DataOffset     uint8
	ACKNumber      uint32
	SequenceNumber uint32
	DstPort        uint16
	SrcPort        uint16
}

func DecodeTCP(raw unsafe.Pointer, length uint32) (*TCPHeader, error) {
	if SizeOfTCPHeader > uintptr(length) {
		return nil, fmt.Errorf("Frame is too short: %d", length)
	}

	return (*TCPHeader)(unsafe.Pointer(uintptr(raw) - SizeOfTCPHeader)), nil
}
