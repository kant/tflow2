package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/taktv6/tflow2/convert"
	"github.com/taktv6/tflow2/database"
)

func (fe *Frontend) prometheusHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	query := database.Query{}
	var errs []error

	// Parse and assign breakdown
	breakdown := strings.Split(params.Get("breakdown"), ",")

	var breakDownError error
	if len(breakdown) == 1 && breakdown[0] == "" {
		breakDownError = errors.New("breakdown parameter missing")
	} else if err := query.Breakdown.Set(breakdown); err != nil {
		breakDownError = fmt.Errorf("breakdown parameter invalid: %s", err)
	}
	if breakDownError != nil {
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "%s\n", breakDownError)
		fmt.Fprintf(buf, "please pass a comma separated list of:\n")
		for _, label := range database.GetBreakdownLabels() {
			fmt.Fprintf(buf, "- %s\n", label)
		}
		errs = append(errs, errors.New(buf.String()))
	}

	// Fetch router parameter
	router := params.Get("router")
	if router == "" {
		errs = append(errs, errors.New("router parameter missing"))
	}

	var ts int64

	// Optional timestamp
	if val := params.Get("ts"); val != "" {
		var err error
		ts, err = strconv.ParseInt(val, 10, 32)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to parse ts: %s", err.Error()))
			return
		}
	} else {
		// Select most recent complete timeslot
		ts = fe.flowDB.CurrentTimeslot() - fe.flowDB.AggregationPeriod()
	}

	if len(errs) > 0 {
		http.Error(w, "Invalid parameters\n", 422)
		for _, err := range errs {
			fmt.Fprintln(w, err.Error())
		}
		return
	}

	// Create conditions for router and timerange
	query.Cond = []database.Condition{
		database.Condition{
			Field:    database.FieldRouter,
			Operator: database.OpEqual,
			Operand:  convert.IPByteSlice(router),
		},
		database.Condition{
			Field:    database.FieldTimestamp,
			Operator: database.OpEqual,
			Operand:  convert.Int64Byte(ts),
		},
	}

	// Run the query
	result, err := fe.flowDB.RunQuery(&query)
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), 502)
		return
	}

	// Empty result?
	if len(result.Timestamps) == 0 {
		http.Error(w, "No data found", 404)
		return
	}

	// Hints for Prometheus
	fmt.Fprintln(w, "# HELP tflow_bytes Bytes transmitted")
	fmt.Fprintln(w, "# TYPE tflow_bytes gauge")

	// Print the data
	for key, val := range result.Data[ts] {
		fmt.Fprintf(w, "tflow_bytes{%s} %d\n", formatBreakdownKey(&key), val)
	}
}

// formats a breakdown key for prometheus
// see tests for examples
func formatBreakdownKey(key *database.BreakdownKey) string {
	result := bytes.Buffer{}

	key.Each(func(key, value string) {
		if result.Len() > 0 {
			result.WriteRune(',')
		}
		result.WriteString(fmt.Sprintf(`%s="%s"`, key, value))
	})

	return result.String()
}
