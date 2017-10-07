// Copyright 2017 EXARING AG. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package sflow

import (
	"fmt"
	"net"
	"testing"
	"unsafe"

	"github.com/taktv6/tflow2/convert"
)

func TestDecode(t *testing.T) {
	s := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 22, 0, 0, 0, 32, 0, 0, 0, 62, 190, 59, 194, 1, 0, 0, 0, 16, 0, 0, 0, 234, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 233, 3, 0, 0, 237, 199, 45, 191, 139, 110, 125, 230, 182, 29, 57, 172, 218, 131, 46, 119, 222, 169, 239, 221, 168, 115, 245, 18, 162, 61, 247, 165, 225, 137, 141, 210, 165, 115, 237, 171, 115, 10, 153, 41, 121, 49, 57, 188, 199, 201, 25, 85, 91, 144, 240, 211, 169, 192, 41, 161, 202, 222, 113, 99, 33, 78, 210, 92, 70, 28, 134, 39, 126, 255, 10, 8, 1, 1, 0, 0, 118, 202, 230, 1, 16, 128, 78, 151, 101, 60, 114, 24, 235, 218, 161, 4, 80, 0, 127, 251, 90, 95, 2, 153, 37, 185, 194, 50, 6, 63, 0, 64, 128, 86, 180, 5, 0, 69, 0, 8, 236, 43, 4, 113, 78, 32, 82, 114, 59, 217, 103, 216, 128, 0, 0, 0, 4, 0, 0, 0, 198, 5, 0, 0, 1, 0, 0, 0, 144, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 190, 2, 0, 0, 168, 2, 0, 0, 0, 0, 0, 0, 64, 127, 94, 90, 224, 3, 0, 0, 144, 2, 0, 0, 197, 164, 97, 81, 232, 0, 0, 0, 1, 0, 0, 0, 22, 0, 0, 0,
		32, 0, 0, 0,
		62, 190, 59, 194, // Next-Hop
		1, 0, 0, 0, // Address Family
		16, 0, 0, 0, // Flow Data Length
		234, 3, 0, 0, // Enterprise/Type (Extended router data)

		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 233, 3, 0, 0, 75, 93, 7, 11, 45, 17, 165, 149, 120, 168, 247, 10, 136, 114, 169, 85, 104, 20, 124, 203, 71, 138, 96, 64, 49, 131, 198, 14, 182, 117, 228, 255, 19, 147, 111, 15, 10, 33, 225, 93, 118, 40, 164, 113, 66, 24, 150, 16, 218, 69, 118, 184, 150, 106, 186, 60, 41, 243, 231, 211, 233, 0, 131, 153, 43, 0, 3, 148, 69, 3, 10, 8, 1, 1, 0, 0, 233, 206, 130, 1, 16, 128, 172, 10, 7, 23, 40, 164, 166, 29, 62, 63, 80, 0, 43, 248, 17, 31, 4, 153, 37, 185, 46, 251, 6, 63, 0, 64, 174, 209, 180, 5, 0, 69, 0, 8, 236, 43, 4, 113, 78, 32, 82, 114, 59, 217, 103, 216, 128, 0, 0, 0, 4, 0, 0, 0, 198, 5, 0, 0, 1, 0, 0, 0, 144, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 190, 2, 0, 0, 170, 2, 0, 0, 0, 0, 0, 0, 96, 123, 94, 90, 224, 3, 0, 0, 144, 2, 0, 0, 196, 164, 97, 81, 232, 0, 0, 0, 1, 0, 0, 0, 14, 0, 0, 0, 32, 0, 0, 0, 57, 96, 89, 195, 1, 0, 0, 0, 16, 0, 0, 0, 234, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 233, 3, 0, 0, 215, 208, 48, 29, 1, 33, 28, 71, 110, 205, 210, 148, 225, 14, 237, 179, 197, 53, 4, 58, 246, 63, 228, 230, 166, 133, 111, 70, 124, 147, 240, 222, 21, 201, 13, 213, 140, 73, 144, 70, 156, 85, 47, 29, 86, 176, 195, 134, 78, 168, 63, 135, 252, 8, 80, 190, 183, 194, 133, 210, 26, 105, 239, 144, 29, 0, 2, 76, 160, 139, 10, 8, 1, 1, 0, 0, 167, 74, 239, 0, 16, 128, 210, 21, 9, 11, 29, 195, 141, 208, 244, 155, 80, 0, 91, 117, 210, 92, 4, 153, 37, 185, 251, 64, 6, 63, 0, 64, 209, 208, 212, 5, 0, 69, 0, 8, 188, 28, 4, 113, 78, 32, 3, 248, 103, 156, 181, 132, 128, 0, 0, 0, 4, 0, 0, 0, 230, 5, 0, 0, 1, 0, 0, 0, 144, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 149, 2, 0, 0, 170, 2, 0, 0, 0, 0, 0, 0, 96, 133, 157, 123, 224, 3, 0, 0, 149, 2, 0, 0, 116, 98, 15, 54, 232, 0, 0, 0, 1, 0, 0, 0, 10, 0, 0, 0, 32, 0, 0, 0, 33, 250, 157, 62, 1, 0, 0, 0, 16, 0, 0, 0, 234, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 233, 3, 0, 0, 193, 111, 105, 60, 190, 220, 121, 229, 158, 159, 65, 27, 79, 59, 89, 152, 153, 147, 249, 41, 34, 174, 115, 106, 7, 8, 148, 19, 165, 47, 135, 86, 42, 17, 129, 84, 254, 130, 222, 106, 42, 106, 209, 185, 205, 208, 71, 17, 126, 140, 32, 197, 254, 206, 15, 11, 174, 65, 151, 178, 9, 214, 21, 70, 123, 1, 217, 142, 46, 12, 10, 8, 1, 1, 0, 0, 80, 121, 23, 4, 16, 128, 116, 173, 164, 116, 56, 194, 157, 44, 176, 189, 80, 0, 246, 113, 186, 87, 3, 153, 37, 185, 75, 197, 6, 63, 0, 64, 255, 84, 212, 5, 0, 69, 0, 8, 185, 28, 4, 113, 78, 32, 148, 2, 127, 31, 113, 128, 128, 0, 0, 0, 4, 0, 0, 0, 230, 5, 0, 0, 1, 0, 0, 0, 144, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 146, 2, 0, 0, 171, 2, 0, 0, 0, 0, 0, 0, 128, 85, 79, 192, 224, 3, 0, 0, 146, 2, 0, 0, 211, 127, 173, 95, 232, 0, 0, 0, 1, 0, 0, 0, 10, 0, 0, 0,
		32, 0, 0, 0,
		33, 250, 157, 62, // Next-Hop
		1, 0, 0, 0, // Address Family
		16, 0, 0, 0, // Flow Data Length
		234, 3, 0, 0, // Enterprise/Type (Extended router data)

		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 210, 0, 0, 0, 16, 0, 0, 0, 233, 3, 0, 0, 209, 50, 196, 16, 191, 134, 236, 166, 206, 27, 249, 140, 64, 231, 148, 246, 19, 88, 36, 9, 167, 240, 97, 133, 46, 175, 100, 47, 143, 160, 84, 35, 234, 71, 176, 116, 103, 119, 151, 133, 184, 52, 169, 202, 53, 231, 149, 40, 16, 81, 31, 242, 100, 122, 152, 78, 32, 133, 116, 22, 89, 122, 149, 27, 64, 0, 173, 248, 203, 199, 10, 8, 1, 1, 0, 0, 199, 212, 235, 0, 16,
		128,               // Header Length
		92, 180, 133, 203, // ACK Number
		31, 4, 191, 24, // Sequence Number
		222, 148, // DST port
		80, 0, // SRC port

		19, 131, 191, 87, // DST IP
		238, 153, 37, 185, // SRC IP
		186, 25, // Header Checksum
		6,     // Protocol
		62,    // TTL
		0, 64, // Flags + Fragment offset
		131, 239, // Identifier
		212, 5, // Total Length
		0,  // TOS
		69, // Version + Length

		0, 8, // EtherType
		185, 28, 4, 113, 78, 32, // Source MAC
		148, 2, 127, 31, 113, 128, // Destination MAC

		128, 0, 0, 0, // Original Packet length
		4, 0, 0, 0, // Payload removed
		230, 5, 0, 0, // Frame length
		1, 0, 0, 0, // Header Protocol
		144, 0, 0, 0, // Flow Data Length
		1, 0, 0, 0, // Enterprise/Type

		3, 0, 0, 0, // Flow Record count
		146, 2, 0, 0, // Output interface
		7, 2, 0, 0, // Input interface
		0, 0, 0, 0, // Dropped Packets
		160, 81, 79, 192, // Sampling Pool
		224, 3, 0, 0, // Sampling Rate
		146, 2, 0, 0, // Source ID + Index
		210, 127, 173, 95, // Sequence Number
		232, 0, 0, 0, // sample length
		1, 0, 0, 0, // Enterprise/Type

		5, 0, 0, 0, // NumSamples
		111, 0, 0, 0, // SysUpTime
		222, 0, 0, 0, // Sequence Number
		0, 0, 0, 0, // Sub-AgentID
		14, 19, 205, 10, // Agent Address
		1, 0, 0, 0, // Agent Address Type
		5, 0, 0, 0, // Version
	}
	s = convert.Reverse(s)

	packet, err := Decode(s, net.IP([]byte{1, 1, 1, 1}))
	if err != nil {
		t.Errorf("Decoding packet failed: %v\n", err)
	}

	if packet.Header.AgentAddress.String() != "10.205.19.14" {
		t.Errorf("Incorrect AgentAddress: Exptected 10.205.19.14 got %s", packet.Header.AgentAddress.String())
	}

	dump(packet)

}

func dump(packet *Packet) {
	fmt.Printf("PACKET DUMP:\n")
	for _, fs := range packet.FlowSamples {
		if fs.ExtendedRouterData != nil {
			fmt.Printf("Extended router data:\n")
			fmt.Printf("Next-Hop: %s\n", fs.ExtendedRouterData.NextHop.String())
		}
		if fs.RawPacketHeader != nil {
			fmt.Printf("Raw packet header:\n")
			fmt.Printf("OriginalPacketLength: %d\n", fs.RawPacketHeader.OriginalPacketLength)
			fmt.Printf("Original Packet:\n")
			data := unsafe.Pointer(uintptr(fs.RawPacketHeader.RawData) - uintptr(fs.RawPacketHeader.OriginalPacketLength))
			ptr := uintptr(data)
			for i := uint32(0); i < fs.RawPacketHeader.OriginalPacketLength; i++ {
				d := (*byte)(unsafe.Pointer(ptr))
				fmt.Printf(" %v ", *d)
				ptr++
			}
			fmt.Printf("\nEND OF DUMP\n")
		}
	}
}

func testEq(a, b []byte) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
