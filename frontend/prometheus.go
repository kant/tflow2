package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/taktv6/tflow2/convert"
	"github.com/taktv6/tflow2/database"
)

const prefix = "tflow_"

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

	// Create a new collector and pass it to Prometheus for handling
	reg := prometheus.NewRegistry()
	reg.MustRegister(newCollector(result, breakdown))

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	}).ServeHTTP(w, r)
}

type collector struct {
	labels    []string
	result    *database.Result
	bytesDesc *prometheus.Desc
}

func newCollector(result *database.Result, labels []string) *collector {
	return &collector{
		result:    result,
		labels:    labels,
		bytesDesc: prometheus.NewDesc(prefix+"bytes", "Bytes transmitted", labels, nil),
	}
}

// Describe writes the descriptions into the channel
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bytesDesc
}

// Collect writes the metrhics into the channel
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	result := c.result

	if len(result.Timestamps) == 0 {
		glog.Errorf("No timestamps found")
		return
	}

	ts := result.Timestamps[0]
	data := result.Data[ts.(int64)]

	for key, val := range data {

		labels := make([]string, len(c.labels))
		for i, label := range c.labels {
			labels[i] = key.Get(label)
		}

		ch <- prometheus.MustNewConstMetric(c.bytesDesc, prometheus.GaugeValue, float64(val), labels...)
	}
}
