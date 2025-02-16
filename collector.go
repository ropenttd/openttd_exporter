package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ropenttd/gopenttd/pkg/gopenttd"
	"time"
)

var (
	// Create a gauge to show whether the server is queryable or not..
	statusDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "active"),
		"Server state.",
		[]string{}, nil,
	)
	// Create a gauge to track user counts. Spectators and overall Clients are
	// differentiated via a "type" label.
	clientsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "clients"),
		"Currently active clients.",
		[]string{"type"}, nil,
	)
	clientsLimitsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "client_limits"),
		"Client limits.",
		[]string{"type"}, nil,
	)

	// Create a gauge to track company count.
	companiesDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "companies"),
		"Currently active companies.",
		[]string{}, nil,
	)
	companiesLimitsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "company_limit"),
		"Company limit.",
		[]string{}, nil,
	)

	// Create a counter for the current runtime.
	gameRunTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "game", "runtime"),
		"How long the current game has been running in in-game days.",
		[]string{}, nil,
	)

	queryTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openttd", "meta", "query_time"),
		"Duration of the last query.",
		[]string{}, nil,
	)
)

type OpenttdCollector struct{}

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc OpenttdCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc OpenttdCollector) Collect(ch chan<- prometheus.Metric) {
	begin := time.Now()
	result, _ := gopenttd.ScanServer(*targetServer, *targetPort)
	duration := time.Since(begin)

	var state int
	if result.Status {
		state = 1
	} else {
		state = 0
	}

	ch <- prometheus.MustNewConstMetric(
		statusDesc,
		prometheus.GaugeValue,
		float64(state),
	)

	ch <- prometheus.MustNewConstMetric(
		clientsDesc,
		prometheus.GaugeValue,
		float64(result.NumClients),
		"clients",
	)
	ch <- prometheus.MustNewConstMetric(
		clientsDesc,
		prometheus.GaugeValue,
		float64(result.NumSpectators),
		"spectators",
	)
	ch <- prometheus.MustNewConstMetric(
		clientsLimitsDesc,
		prometheus.GaugeValue,
		float64(result.MaxClients),
		"clients",
	)
	ch <- prometheus.MustNewConstMetric(
		clientsLimitsDesc,
		prometheus.GaugeValue,
		float64(result.MaxSpectators),
		"spectators",
	)

	ch <- prometheus.MustNewConstMetric(
		companiesDesc,
		prometheus.GaugeValue,
		float64(result.NumCompanies),
	)
	ch <- prometheus.MustNewConstMetric(
		companiesLimitsDesc,
		prometheus.GaugeValue,
		float64(result.MaxCompanies),
	)

	ch <- prometheus.MustNewConstMetric(
		gameRunTimeDesc,
		prometheus.CounterValue,
		float64(result.DateCurrent.Sub(result.DateStart).Hours()/24),
	)

	ch <- prometheus.MustNewConstMetric(
		queryTimeDesc,
		prometheus.GaugeValue,
		duration.Seconds(),
	)

}

func NewOpenttdCollector(reg prometheus.Registerer) *OpenttdCollector {
	cc := &OpenttdCollector{}
	prometheus.WrapRegistererWith(prometheus.Labels{}, reg).MustRegister(cc)
	return cc
}
