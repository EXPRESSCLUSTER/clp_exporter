package collector

import (
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "clp"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "controller_duration_seconds"),
		"sample_exporter: Duration of a collector scrape",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "controller_success"),
		"sample_exporter: Whether a collector succeeded",
		[]string{"collector"},
		nil,
	)
	factories      = make(map[string]func() (Collector, error))
	collectorState = make(map[string]int)
)

func registCollector(collector string, f func() (Collector, error)) {
	factories[collector] = f
	collectorState[collector] = 0
}

type CLPCollector struct {
	Collectors map[string]Collector
}

type Collector interface {
	Update(ch chan<- prometheus.Metric) error
}

func NewCLPCollector() (*CLPCollector, error) {
	collectors := make(map[string]Collector)
	for k := range collectorState {
		f, err := factories[k]()
		if err != nil {
			return nil, err
		}
		collectors[k] = f
	}
	return &CLPCollector{Collectors: collectors}, nil
}

func (sc CLPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func (sc CLPCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(sc.Collectors))
	for name, c := range sc.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Printf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err.Error())
		success = 0
	}
	success = 1
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}
