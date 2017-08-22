package frontend

import (
	"fmt"
	"net"
	"strconv"

	"github.com/taktv6/tflow2/convert"
	"github.com/taktv6/tflow2/database"
)

// ConditionsExt is external representation of conditions of a query
type ConditionsExt []ConditionExt

// ConditionExt is external representation of a query condition
type ConditionExt struct {
	Field    string
	Operator int
	Operand  string
}

// QueryExt represents a query in the way it is received from the frontend
type QueryExt struct {
	Cond      ConditionsExt
	Breakdown database.BreakdownFlags
	TopN      int
}

func (ext *ConditionExt) toQueryExt() (*database.Condition, error) {
	var operand []byte
	fieldNum := database.GetFieldByName(ext.Field)

	switch fieldNum {
	case database.FieldTimestamp:
		op, err := strconv.Atoi(ext.Operand)
		if err != nil {
			return nil, err
		}
		operand = convert.Int64Byte(int64(op))

	case database.FieldProtocol, database.FieldSrcPort, database.FieldDstPort, database.FieldIntIn, database.FieldIntOut:
		op, err := strconv.Atoi(ext.Operand)
		if err != nil {
			return nil, err
		}
		operand = convert.Uint16Byte(uint16(op))

	case database.FieldSrcAddr, database.FieldDstAddr, database.FieldRouter, database.FieldNextHop:
		operand = convert.IPByteSlice(ext.Operand)

	case database.FieldSrcAs, database.FieldDstAs, database.FieldNextHopAs:
		op, err := strconv.Atoi(ext.Operand)
		if err != nil {
			return nil, err
		}
		operand = convert.Uint32Byte(uint32(op))

	case database.FieldSrcPfx, database.FieldDstPfx:
		_, pfx, err := net.ParseCIDR(string(ext.Operand))
		if err != nil {
			return nil, err
		}
		operand = []byte(pfx.String())
	default:
		return nil, fmt.Errorf("unknown field: %s", ext.Field)
	}

	return &database.Condition{
		Field:    fieldNum,
		Operator: ext.Operator,
		Operand:  operand,
	}, nil
}

// translateQuery translates a query from external representation to internal representaion
func translateQuery(e *QueryExt) (*database.Query, error) {
	var q database.Query
	q.Breakdown = e.Breakdown
	q.TopN = e.TopN

	for _, c := range e.Cond {
		cond, err := c.toQueryExt()

		if err != nil {
			return nil, err
		}
		q.Cond = append(q.Cond, *cond)
	}

	return &q, nil
}
