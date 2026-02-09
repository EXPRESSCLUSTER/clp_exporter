package collector

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	monitorSubsystem = "monitor"
)

type monitorCollector struct{}

func init() {
	registCollector(monitorSubsystem, NewMonitorCollector)
}

func NewMonitorCollector() (Collector, error) {
	return &monitorCollector{}, nil
}

func (c *monitorCollector) Update(ch chan<- prometheus.Metric) error {
	metricType := prometheus.GaugeValue

	// Get all monitor types from clp.conf file
	clpcfget_montypes := exec.Command("clpcfget", "-e", "/root/monitor/types")
	output, err := clpcfget_montypes.Output()
	if err != nil {
		log.Fatalf("Failed to execute clpcfget -e /root/monitor/types: %v", err)
	}
	log.Printf("Types:\n%s", strings.TrimSpace(string(output)))

	// Get monitor resource name and run clpperfc command
	montypes := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, montype := range montypes {
		if strings.TrimSpace(montype) != "" {
			log.Printf("Type: %s", montype)
			clpcfget_montype := exec.Command("clpcfget", "-e", "/root/monitor/"+montype)
			output, err = clpcfget_montype.Output()
			if err != nil {
				log.Fatalf("Failed to execute clpcfget -e /root/monitor/%s: %v", montype, err)
				continue
			}
			monitors := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, monitor := range monitors {
				if strings.TrimSpace(monitor) != "" {
					log.Printf("Name: %s", monitor)
				}
				// Run clpperfc -m <monitor> command
				clpperfc_mon := exec.Command("clpperfc", "-m", monitor)
				output, err = clpperfc_mon.Output()
				if err != nil {
					log.Fatalf("Failed to execute clpperfc -m %s: %v", monitor, err)
				}
				log.Printf("clpperfc -m %s: %s", monitor, strings.TrimSpace(string(output)))
				value, err := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
				if err != nil {
					log.Fatalf("Failed to parse output of clpperfc -m %s: %v", monitor, err)
				}
				log.Printf("Value: %d", value)
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, monitorSubsystem, monitor),
						fmt.Sprintf("Monitor %s of type %s", monitor, montype),
						nil, nil,
					),
					metricType, float64(value),
				)
			}
		}
	}

	return nil
}
