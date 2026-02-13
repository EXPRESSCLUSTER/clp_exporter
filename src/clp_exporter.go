package main

import (
	"clp_exporter/collector"
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen address", ":29090", "The Address to listen on for HTTP Requests.")
)

func main() {
	flag.Parse()

	c, err := collector.NewCLPCollector()
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
