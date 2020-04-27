package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	addr         = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	targetServer = flag.String("target.host", "localhost", "The host to monitor.")
	targetPort   = flag.Int("target.port", 3979, "The port to monitor.")
)

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		return
	})
}

func main() {
	flag.Parse()

	reg := prometheus.NewPedanticRegistry()

	NewOpenttdCollector(reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	http.Handle("/health", healthz())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
