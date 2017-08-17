package frontend

import (
	"fmt"
	"net/http"
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

	// Parse and assign breakdown
	breakdown := strings.Split(params.Get("breakdown"), ",")
	if len(breakdown) == 0 {
		http.Error(w, "breakdown parameter missing", 422)
		return
	}
	if err := query.Breakdown.Set(breakdown); err != nil {
		http.Error(w, fmt.Sprintf("breakdown parameter invalid: %s", err), 422)
		return
	}

	// Fetch router parameter
	router := params.Get("router")
	if router == "" {
		http.Error(w, "router parameter missing", 422)
		return
	}

	ts := fe.flowDB.CurrentTimeslot() - fe.flowDB.AggregationPeriod()

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
