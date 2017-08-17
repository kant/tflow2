package database

import (
	"net"

	"github.com/golang/glog"
	"github.com/taktv6/tflow2/avltree"
	"github.com/taktv6/tflow2/convert"
)

// TimeGroup groups all indices to flows of a particular router at a particular
// time into one object
type TimeGroup struct {
	Any       *mapTree // Workaround: Why a map? Because: cannot assign to flows[fl.Timestamp][rtr].Any
	SrcAddr   *mapTree
	DstAddr   *mapTree
	Protocol  *mapTree
	IntIn     *mapTree
	IntOut    *mapTree
	NextHop   *mapTree
	SrcAs     *mapTree
	DstAs     *mapTree
	NextHopAs *mapTree
	SrcPfx    *mapTree
	DstPfx    *mapTree
	SrcPort   *mapTree
	DstPort   *mapTree
}

func (tg *TimeGroup) filterAndBreakdown(resSum *concurrentResSum, q *Query) BreakdownMap {
	// candidates keeps a list of all trees that fulfill the queries criteria
	candidates := make([]*avltree.Tree, 0)
	for _, c := range q.Cond {
		switch c.Field {
		case FieldTimestamp:
			continue
		case FieldRouter:
			continue
		case FieldProtocol:
			candidates = append(candidates, tg.Protocol.Get(c.Operand[0]))
		case FieldSrcAddr:
			candidates = append(candidates, tg.SrcAddr.Get(net.IP(c.Operand)))
		case FieldDstAddr:
			candidates = append(candidates, tg.DstAddr.Get(net.IP(c.Operand)))
		case FieldIntIn:
			candidates = append(candidates, tg.IntIn.Get(convert.Uint16b(c.Operand)))
		case FieldIntOut:
			candidates = append(candidates, tg.IntOut.Get(convert.Uint16b(c.Operand)))
		case FieldNextHop:
			candidates = append(candidates, tg.NextHop.Get(net.IP(c.Operand)))
		case FieldSrcAs:
			candidates = append(candidates, tg.SrcAs.Get(convert.Uint32b(c.Operand)))
		case FieldDstAs:
			candidates = append(candidates, tg.DstAs.Get(convert.Uint32b(c.Operand)))
		case FieldNextHopAs:
			candidates = append(candidates, tg.NextHopAs.Get(convert.Uint32b(c.Operand)))
		case FieldSrcPort:
			candidates = append(candidates, tg.SrcPort.Get(c.Operand))
		case FieldDstPort:
			candidates = append(candidates, tg.DstPort.Get(c.Operand))
		case FieldSrcPfx:
			candidates = append(candidates, tg.SrcPfx.Get(c.Operand))
		case FieldDstPfx:
			candidates = append(candidates, tg.DstPfx.Get(c.Operand))
		}
	}

	if len(candidates) == 0 {
		candidates = append(candidates, tg.Any.Get(anyIndex))
	}

	// Find common elements of candidate trees
	res := avltree.Intersection(candidates)
	if res == nil {
		glog.Warningf("Interseciton Result was empty!")
		res = tg.Any.Get(anyIndex)
	}

	// Breakdown
	resTime := make(BreakdownMap)
	res.Each(breakdown, q.Breakdown, resSum, resTime)
	return resTime
}
