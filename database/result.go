package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// Result is the result of a query
type Result struct {
	TopKeys     map[BreakdownKey]void
	Timestamps  []interface{}          // sorted timestamps
	Data        map[int64]BreakdownMap // timestamps -> keys -> values
	Aggregation int64
	SampleRate  int
}

// WriteCSV writes the result as CSV into the writer
func (res *Result) WriteCSV(writer io.Writer) {
	w := csv.NewWriter(writer)
	defer w.Flush()

	// Construct table header
	headLine := make([]string, 0)
	headLine = append(headLine, "Time")
	topKeys := make([]BreakdownKey, 0)

	for k := range res.TopKeys {
		topKeys = append(topKeys, k)
		headLine = append(headLine, k.Join("%s:%s"))
	}
	headLine = append(headLine, "Rest")
	w.Write(headLine)

	for _, ts := range res.Timestamps {
		line := make([]string, 0)
		t := time.Unix(ts.(int64), 0)
		line = append(line, fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()))

		// Top flows
		buckets := res.Data[ts.(int64)]
		for _, k := range topKeys {
			if _, ok := buckets[k]; !ok {
				line = append(line, "0")
			} else {
				line = append(line, fmt.Sprintf("%d", buckets[k]/uint64(res.Aggregation)*8*uint64(res.SampleRate)))
			}
		}

		// Remaining flows
		var rest uint64
		for k, v := range buckets {
			if _, ok := res.TopKeys[k]; !ok {
				rest += v
			}
		}
		w.Write(append(line, fmt.Sprintf("%d", rest)))
	}
}
