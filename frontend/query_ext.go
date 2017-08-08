package frontend

import (
	"net"
	"strconv"

	"github.com/taktv6/tflow2/convert"
	"github.com/taktv6/tflow2/database"
)

// ConditionsExt is external representation of conditions of a query
type ConditionsExt []ConditionExt

// ConditionExt is external representation of a query condition
type ConditionExt struct {
	Field    int
	Operator int
	Operand  string
}

// QueryExt represents a query in the way it is received from the frontend
type QueryExt struct {
	Cond      ConditionsExt
	Breakdown database.BreakdownFlags
	TopN      int
}

// translateQuery translates a query from external representation to internal representaion
func translateQuery(e *QueryExt) (*database.Query, error) {
	var q database.Query
	q.Breakdown = e.Breakdown
	q.TopN = e.TopN

	for _, c := range e.Cond {
		var operand []byte

		switch c.Field {
		case database.FieldTimestamp:
			op, err := strconv.Atoi(c.Operand)
			if err != nil {
				return nil, err
			}
			operand = convert.Int64Byte(int64(op))

		case database.FieldProtocol, database.FieldSrcPort, database.FieldDstPort, database.FieldIntIn, database.FieldIntOut:
			op, err := strconv.Atoi(c.Operand)
			if err != nil {
				return nil, err
			}
			operand = convert.Uint16Byte(uint16(op))

		case database.FieldSrcAddr, database.FieldDstAddr, database.FieldRouter, database.FieldNextHop:
			operand = convert.IPByteSlice(c.Operand)

		case database.FieldSrcAs, database.FieldDstAs, database.FieldNextHopAs:
			op, err := strconv.Atoi(c.Operand)
			if err != nil {
				return nil, err
			}
			operand = convert.Uint32Byte(uint32(op))

		case database.FieldSrcPfx, database.FieldDstPfx:
			_, pfx, err := net.ParseCIDR(string(c.Operand))
			if err != nil {
				return nil, err
			}
			operand = []byte(pfx.String())
		}

		q.Cond = append(q.Cond, database.Condition{
			Field:    c.Field,
			Operator: c.Operator,
			Operand:  operand,
		})
	}

	return &q, nil
}
