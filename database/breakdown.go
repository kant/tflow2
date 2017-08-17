package database

import (
	"bytes"
	"fmt"
	"net"

	"github.com/golang/glog"
	"github.com/taktv6/tflow2/avltree"
	"github.com/taktv6/tflow2/netflow"
)

// BreakdownKey is the key used for the brakedown map
type BreakdownKey [FieldMax]string

// BreakdownMap maps breakdown keys to values
type BreakdownMap map[BreakdownKey]uint64

// BreakdownFlags defines by what fields data should be broken down in a query
type BreakdownFlags struct {
	Router     bool
	Family     bool
	SrcAddr    bool
	DstAddr    bool
	Protocol   bool
	IntIn      bool
	IntOut     bool
	NextHop    bool
	SrcAsn     bool
	DstAsn     bool
	NextHopAsn bool
	SrcPfx     bool
	DstPfx     bool
	SrcPort    bool
	DstPort    bool
}

var breakdownLabels = map[int]string{
	FieldSrcAddr:   "Src",
	FieldDstAddr:   "Dst",
	FieldProtocol:  "Proto",
	FieldIntIn:     "IntIn",
	FieldIntOut:    "IntOut",
	FieldNextHop:   "NH",
	FieldSrcAs:     "SrcAsn",
	FieldDstAs:     "DstAsn",
	FieldNextHopAs: "NH_AS",
	FieldSrcPfx:    "SrcNet",
	FieldDstPfx:    "DstNet",
	FieldSrcPort:   "SrcPort",
	FieldDstPort:   "DstPort",
}

// reverse mapping for breakdownLabels
func breakdownIndex(key string) int {
	for i, k := range breakdownLabels {
		if k == key {
			return i
		}
	}
	return -1
}

// Set sets the value of a field
func (bk *BreakdownKey) Set(key string, value string) {
	bk[breakdownIndex(key)] = value
}

// Get returns the value of a field
func (bk *BreakdownKey) Get(key string) string {
	return bk[breakdownIndex(key)]
}

// String builds a textual representation of the key
func (bk *BreakdownKey) String() string {
	var buffer bytes.Buffer

	for i, val := range bk {
		if val == "" {
			continue
		}
		if label, ok := breakdownLabels[i]; ok {
			if buffer.Len() > 0 {
				buffer.WriteRune(',')
			}
			buffer.WriteString(label + ":" + val)
		}
	}

	return buffer.String()
}

// Set enables the flags in the given list
func (bf *BreakdownFlags) Set(keys []string) error {
	for _, key := range keys {
		switch key {
		case "Router":
			bf.Router = true
		case "Family":
			bf.Family = true
		case "SrcAddr":
			bf.SrcAddr = true
		case "DstAddr":
			bf.DstAddr = true
		case "Protocol":
			bf.Protocol = true
		case "IntIn":
			bf.IntIn = true
		case "IntOut":
			bf.IntOut = true
		case "NextHop":
			bf.NextHop = true
		case "SrcAsn":
			bf.SrcAsn = true
		case "DstAsn":
			bf.DstAsn = true
		case "NextHopAsn":
			bf.NextHopAsn = true
		case "SrcPfx":
			bf.SrcPfx = true
		case "DstPfx":
			bf.DstPfx = true
		case "SrcPort":
			bf.SrcPort = true
		case "DstPort":
			bf.DstPort = true
		default:
			return fmt.Errorf("invalid breakdown key: %s", key)
		}
	}
	return nil
}

// breakdown build all possible relevant keys of flows for flows in tree `node`
// and builds sums for each key in order to allow us to find top combinations
func breakdown(node *avltree.TreeNode, vals ...interface{}) {
	if len(vals) != 3 {
		glog.Errorf("lacking arguments")
		return
	}

	bd := vals[0].(BreakdownFlags)
	sums := vals[1].(*concurrentResSum)
	buckets := vals[2].(BreakdownMap)
	fl := node.Value.(*netflow.Flow)

	key := BreakdownKey{}

	if bd.SrcAddr {
		key[FieldSrcAddr] = net.IP(fl.SrcAddr).String()
	}
	if bd.DstAddr {
		key[FieldDstAddr] = net.IP(fl.DstAddr).String()
	}
	if bd.Protocol {
		key[FieldProtocol] = fmt.Sprintf("%d", fl.Protocol)
	}
	if bd.IntIn {
		key[FieldIntIn] = fmt.Sprintf("%d", fl.IntIn)
	}
	if bd.IntOut {
		key[FieldIntOut] = fmt.Sprintf("%d", fl.IntOut)
	}
	if bd.NextHop {
		key[FieldNextHop] = net.IP(fl.NextHop).String()
	}
	if bd.SrcAsn {
		key[FieldSrcAs] = fmt.Sprintf("%d", fl.SrcAs)
	}
	if bd.DstAsn {
		key[FieldDstAs] = fmt.Sprintf("%d", fl.DstAs)
	}
	if bd.NextHopAsn {
		key[FieldNextHopAs] = fmt.Sprintf("%d", fl.NextHopAs)
	}
	if bd.SrcPfx {
		if fl.SrcPfx != nil {
			key[FieldSrcPfx] = fl.SrcPfx.ToIPNet().String()
		} else {
			key[FieldSrcPfx] = "0.0.0.0/0"
		}
	}
	if bd.DstPfx {
		if fl.DstPfx != nil {
			key[FieldDstPfx] = fl.DstPfx.ToIPNet().String()
		} else {
			key[FieldDstPfx] = "0.0.0.0/0"
		}
	}
	if bd.SrcPort {
		key[FieldSrcPort] = fmt.Sprintf("%d", fl.SrcPort)
	}
	if bd.DstPort {
		key[FieldDstPort] = fmt.Sprintf("%d", fl.DstPort)
	}

	// Build sum for key
	buckets[key] += fl.Size

	// Build overall sum
	sums.Lock.Lock()
	sums.Values[key] += fl.Size
	sums.Lock.Unlock()
}
