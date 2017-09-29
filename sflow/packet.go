package sflow

import (
	"net"
	"unsafe"
)

// Packet is a decoded representation of a single sflow UDP packet.
type Packet struct {
	// A pointer to the packets headers
	Header *Header

	headerTop    *headerTop
	headerBottom *headerBottom

	FlowSamples []*FlowSample

	// A slice of pointers to FlowSet. Each element is instance of (Data)FlowSet
	// found in this packet
	//FlowSamples []*FlowSample

	// Buffer is a slice pointing to the original byte array that this packet was decoded from.
	// This field is only populated if debug level is at least 2
	Buffer []byte
}

var (
	sizeOfHeaderTop                = unsafe.Sizeof(headerTop{})
	sizeOfHeaderBottom             = unsafe.Sizeof(headerBottom{})
	sizeOfFlowSampleHeader         = unsafe.Sizeof(FlowSampleHeader{})
	sizeOfRawPacketHeader          = unsafe.Sizeof(RawPacketHeader{})
	sizeofExtendedRouterData       = unsafe.Sizeof(ExtendedRouterData{})
	sizeOfextendedRouterDataTop    = unsafe.Sizeof(extendedRouterDataTop{})
	sizeOfextendedRouterDataBottom = unsafe.Sizeof(extendedRouterDataBottom{})
)

type Header struct {
	Version          uint32
	AgentAddressType uint32
	AgentAddress     net.IP
	SubAgentID       uint32
	SequenceNumber   uint32
	SysUpTime        uint32
	NumSamples       uint32
}

type headerTop struct {
	AgentAddressType uint32
	Version          uint32
}

type headerBottom struct {
	NumSamples     uint32
	SysUpTime      uint32
	SequenceNumber uint32
	SubAgentID     uint32
}

type FlowSample struct {
	FlowSampleHeader    *FlowSampleHeader
	RawPacketHeader     *RawPacketHeader
	RawPacketHeaderData unsafe.Pointer
	ExtendedRouterData  *ExtendedRouterData
}

type FlowSampleHeader struct {
	FlowRecord         uint32
	OutputIf           uint32
	InputIf            uint32
	DroppedPackets     uint32
	SamplePool         uint32
	SamplingRate       uint32
	SourceIDClassIndex uint32
	SequenceNumber     uint32
	SampleLength       uint32
	EnterpriseType     uint32
}

type RawPacketHeader struct {
	OriginalPacketLength uint32
	PayloadRemoved       uint32
	FrameLength          uint32
	HeaderProtocol       uint32
	FlowDataLength       uint32
	EnterpriseType       uint32
}

type extendedRouterDataTop struct {
	AddressType    uint32
	FlowDataLength uint32
	EnterpriseType uint32
}

type extendedRouterDataBottom struct {
	NextHopDestinationMask uint32
	NextHopSourceMask      uint32
}

type ExtendedRouterData struct {
	NextHopDestinationMask uint32
	NextHopSourceMask      uint32
	NextHop                net.IP
	AddressType            uint32
	FlowDataLength         uint32
	EnterpriseType         uint32
}
