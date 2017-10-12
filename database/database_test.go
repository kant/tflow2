package database

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/taktv6/tflow2/convert"
	"github.com/taktv6/tflow2/iana"
	"github.com/taktv6/tflow2/intfmapper"
	"github.com/taktv6/tflow2/netflow"
)

type intfMapper struct {
}

func (m *intfMapper) GetInterfaceIDByName(agent string) intfmapper.InterfaceIDByName {
	return intfmapper.InterfaceIDByName{
		"xe-0/0/1": 1,
		"xe-0/0/2": 2,
		"xe-0/0/3": 3,
	}
}

func (m *intfMapper) GetInterfaceNameByID(agent string) intfmapper.InterfaceNameByID {
	return intfmapper.InterfaceNameByID{
		1: "xe-0/0/1",
		2: "xe-0/0/2",
		3: "xe-0/0/3",
	}
}

func TestQuery(t *testing.T) {
	minute := int64(60)
	hour := int64(3600)

	ts1 := int64(3600)
	ts1 = ts1 - ts1%minute

	tests := []struct {
		name           string
		flows          []*netflow.Flow
		query          *Query
		expectedResult Result
	}{
		{
			/*
				Testcase: 2 flows from AS100 to AS300 and back (TCP session).
			*/
			name: "Test 1",
			flows: []*netflow.Flow{
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 1},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1000,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 2},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1000,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 1},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 2},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
			},
			query: &Query{
				Cond: []Condition{
					{
						Field:    FieldAgent,
						Operator: OpEqual,
						Operand:  []byte("test01.pop01"),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpGreater,
						Operand:  convert.Uint64Byte(uint64(ts1 - 3*minute)),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpSmaller,
						Operand:  convert.Uint64Byte(uint64(ts1 + minute)),
					},
					{
						Field:    FieldIntOut,
						Operator: OpEqual,
						Operand:  convert.Uint16Byte(uint16(1)),
					},
				},
				Breakdown: BreakdownFlags{
					SrcAddr: true,
					DstAddr: true,
				},
				TopN: 100,
			},
			expectedResult: Result{
				TopKeys: map[BreakdownKey]void{
					BreakdownKey{
						FieldSrcAddr: "30.0.0.1",
						FieldDstAddr: "10.0.0.1",
					}: void{},
					BreakdownKey{
						FieldSrcAddr: "30.0.0.2",
						FieldDstAddr: "10.0.0.1",
					}: void{},
				},
				Timestamps: []int64{
					ts1,
				},
				Data: map[int64]BreakdownMap{
					ts1: BreakdownMap{
						BreakdownKey{
							FieldSrcAddr: "30.0.0.1",
							FieldDstAddr: "10.0.0.1",
						}: 40000,
						BreakdownKey{
							FieldSrcAddr: "30.0.0.2",
							FieldDstAddr: "10.0.0.1",
						}: 40000,
					},
				},
				Aggregation: minute,
			},
		},

		{
			/*
				Testcase: 2 flows from AS100 to AS300 and back (TCP session).
				Opposite direction of Test 1
			*/
			name: "Test 2",
			flows: []*netflow.Flow{
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 1},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1000,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 2},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1000,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 1},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 2},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
			},
			query: &Query{
				Cond: []Condition{
					{
						Field:    FieldAgent,
						Operator: OpEqual,
						Operand:  []byte("test01.pop01"),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpGreater,
						Operand:  convert.Uint64Byte(uint64(ts1 - 3*minute)),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpSmaller,
						Operand:  convert.Uint64Byte(uint64(ts1 + minute)),
					},
					{
						Field:    FieldIntOut,
						Operator: OpEqual,
						Operand:  convert.Uint16Byte(uint16(3)),
					},
				},
				Breakdown: BreakdownFlags{
					SrcAddr: true,
					DstAddr: true,
				},
				TopN: 100,
			},
			expectedResult: Result{
				TopKeys: map[BreakdownKey]void{
					BreakdownKey{
						FieldSrcAddr: "10.0.0.1",
						FieldDstAddr: "30.0.0.1",
					}: void{},
					BreakdownKey{
						FieldSrcAddr: "10.0.0.1",
						FieldDstAddr: "30.0.0.2",
					}: void{},
				},
				Timestamps: []int64{
					ts1,
				},
				Data: map[int64]BreakdownMap{
					ts1: BreakdownMap{
						BreakdownKey{
							FieldSrcAddr: "10.0.0.1",
							FieldDstAddr: "30.0.0.1",
						}: 4000,
						BreakdownKey{
							FieldSrcAddr: "10.0.0.1",
							FieldDstAddr: "30.0.0.2",
						}: 4000,
					},
				},
				Aggregation: minute,
			},
		},

		{
			/*
				Testcase: 2 flows from AS100 to AS300 and back (TCP session).
				Test TopN function
			*/
			name: "Test 3",
			flows: []*netflow.Flow{
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 1},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1001,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{10, 0, 0, 1},
					DstAddr:    []byte{30, 0, 0, 2},
					Protocol:   6,
					SrcPort:    12345,
					DstPort:    443,
					Packets:    2,
					Size:       1000,
					IntIn:      1,
					IntOut:     3,
					NextHop:    []byte{30, 0, 0, 100},
					SrcAs:      100,
					DstAs:      300,
					NextHopAs:  300,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 1},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
				&netflow.Flow{
					Router:     []byte{1, 2, 3, 4},
					Family:     4,
					SrcAddr:    []byte{30, 0, 0, 2},
					DstAddr:    []byte{10, 0, 0, 1},
					Protocol:   6,
					SrcPort:    443,
					DstPort:    12345,
					Packets:    5,
					Size:       10000,
					IntIn:      3,
					IntOut:     1,
					NextHop:    []byte{10, 0, 0, 100},
					SrcAs:      300,
					DstAs:      100,
					NextHopAs:  100,
					Samplerate: 4,
					Timestamp:  ts1,
				},
			},
			query: &Query{
				Cond: []Condition{
					{
						Field:    FieldAgent,
						Operator: OpEqual,
						Operand:  []byte("test01.pop01"),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpGreater,
						Operand:  convert.Uint64Byte(uint64(ts1 - 3*minute)),
					},
					{
						Field:    FieldTimestamp,
						Operator: OpSmaller,
						Operand:  convert.Uint64Byte(uint64(ts1 + minute)),
					},
					{
						Field:    FieldIntOut,
						Operator: OpEqual,
						Operand:  convert.Uint16Byte(uint16(3)),
					},
				},
				Breakdown: BreakdownFlags{
					SrcAddr: true,
					DstAddr: true,
				},
				TopN: 1,
			},
			expectedResult: Result{
				TopKeys: map[BreakdownKey]void{
					BreakdownKey{
						FieldSrcAddr: "10.0.0.1",
						FieldDstAddr: "30.0.0.1",
					}: void{},
				},
				Timestamps: []int64{
					ts1,
				},
				Data: map[int64]BreakdownMap{
					ts1: BreakdownMap{
						BreakdownKey{
							FieldSrcAddr: "10.0.0.1",
							FieldDstAddr: "30.0.0.1",
						}: 4004,
					},
				},
				Aggregation: minute,
			},
		},
	}

	for _, test := range tests {
		fdb := New(minute, hour, 1, 0, 6, nil, false, &intfMapper{}, map[string]string{
			net.IP([]byte{1, 2, 3, 4}).String(): "test01.pop01",
		}, iana.New())

		for _, flow := range test.flows {
			fdb.Input <- flow
		}

		time.Sleep(time.Second)

		result, err := fdb.RunQuery(test.query)
		if err != nil {
			t.Errorf("Unexpected error on RunQuery: %v", err)
		}

		for k := range result.TopKeys {
			for ts := range result.Data {
				x := result.Data[ts]
				v := x[k]
				fmt.Printf("TS: %d, Key: %v, Value: %d\n", ts, k, v)
			}
		}

		if err := compareResults(test.expectedResult, *result); err != nil {
			t.Errorf("%v", err)
		}

	}
}

func compareResults(res1 Result, res2 Result) error {
	if err := compareTopKeys(res1.TopKeys, res2.TopKeys); err != nil {
		return err
	}

	if err := compareTimestamps(res1.Timestamps, res2.Timestamps); err != nil {
		return err
	}

	if res1.Aggregation != res2.Aggregation {
		return fmt.Errorf("Aggregation %d != %d", res1.Aggregation, res2.Aggregation)
	}

	if err := compareData(res1.Data, res2.Data); err != nil {
		return err
	}

	return nil
}

func compareData(data1 map[int64]BreakdownMap, data2 map[int64]BreakdownMap) error {
	for ts, bdm := range data1 {
		if _, ok := data2[ts]; !ok {
			return fmt.Errorf("TS %d does not exist in data2", ts)
		}

		for bdk := range bdm {
			if _, ok := data2[ts][bdk]; !ok {
				return fmt.Errorf("BreakDownKey %v exists in data1 but not in data2 for TS=%d", bdk, ts)
			}

			if data1[ts][bdk] != data2[ts][bdk] {
				return fmt.Errorf("Values for TS=%d BDK=%v differing: %d != %d", ts, bdk, data1[ts][bdk], data2[ts][bdk])
			}
		}
	}

	return nil
}

func compareTimestamps(ts1 []int64, ts2 []int64) error {
	for i, a := range ts1 {
		if ts2[i] != a {
			return fmt.Errorf("TS %d != %d", a, ts2[i])
		}
	}

	return nil
}

func compareTopKeys(tk1 map[BreakdownKey]void, tk2 map[BreakdownKey]void) error {
	for k1 := range tk1 {
		if _, ok := tk2[k1]; !ok {
			return fmt.Errorf("Top key %v is in tk1 but not in tk2", k1)
		}
	}

	for k2 := range tk2 {
		if _, ok := tk1[k2]; !ok {
			return fmt.Errorf("Top key %v is in tk2 but not in tk1", k2)
		}
	}

	return nil
}
