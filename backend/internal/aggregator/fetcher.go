package aggregator

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type Fetcher struct {
	metricsURL string
	client     *http.Client
}

func NewFetcher(metricsURL string) *Fetcher {
	return &Fetcher{
		metricsURL: metricsURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (f *Fetcher) FetchOSPF() (*frrProto.OSPFMetrics, error) {
	rawData, err := f.fetchRawMetrics()
	if err != nil {
		return nil, err
	}

	return parseOSPFMetrics(rawData)
}

func (f *Fetcher) fetchRawMetrics() ([]byte, error) {
	resp, err := f.client.Get(f.metricsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func parseOSPFMetrics(rawData []byte) (*frrProto.OSPFMetrics, error) {
	var metrics frrProto.OSPFMetrics
	//var OSPFNeighbor []*frrProto.OSPFNeighbor

	parser := expfmt.TextParser{}
	parsedMetrics, err := parser.TextToMetricFamilies(strings.NewReader(string(rawData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	OSPFNeighbors := getNeighbors(parsedMetrics)
	OSPFRoutes := getRoutes(parsedMetrics)
	OSPFInterfaces := getInterfaces(parsedMetrics)
	OSPFlsa := getLsas(parsedMetrics)

	metrics.Neighbors = OSPFNeighbors
	metrics.Routes = OSPFRoutes
	metrics.Interfaces = OSPFInterfaces
	metrics.Lsas = OSPFlsa

	//fmt.Printf("Start: Resulting metrics\n")
	//fmt.Printf("metrics neighbors: \n%v\n\n", metrics.Neighbors)
	//fmt.Printf("metrics routes: \n%v\n\n", metrics.Routes)
	//fmt.Printf("metrics interfaces: \n%v\n\n", metrics.Interfaces)
	//fmt.Printf("metrics lsas: \n%v\n\n", metrics.Lsas)
	//fmt.Printf("End: Resulting metrics\n")
	//if err := json.Unmarshal(rawData, &metrics); err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	//}

	return &metrics, nil
}

func removeComments(metrics string) []byte {
	lines := strings.Split(metrics, "\n")
	var cleaned []string
	for _, line := range lines {
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		cleaned = append(cleaned, line)
	}
	fmt.Println(cleaned)
	return []byte(strings.Join(cleaned, "\n"))
}

func getNeighbors(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFNeighbor {
	var neighbors []*frrProto.OSPFNeighbor

	//fmt.Println(metrics["frr_ospf_neighbor_state"].GetMetric())
	for _, value := range metrics["frr_ospf_neighbor_state"].GetMetric() {
		var neighbor frrProto.OSPFNeighbor

		neighbor.Area = value.GetLabel()[0].GetValue()
		neighbor.Interface = value.GetLabel()[1].GetValue()
		neighbor.Id = value.GetLabel()[2].GetValue()
		neighbor.Ip = value.GetLabel()[3].GetValue()
		neighbors = append(neighbors, &neighbor)
	}

	return neighbors
}

func getInterfaces(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFInterface {
	var result []*frrProto.OSPFInterface
	/*
			  Name string
		    Area string

		    NbrCount int32
		    NbrAdj int32
		    Passive bool
	*/
	//fmt.Println(metrics["frr_ospf_neighbor_state"].GetMetric()[1].GetLabel()[1])
	fmt.Println(metrics["frr_ospf_neighbor_state"])
	//for _, value := range metrics["frr_ospf_neighbor_state"].GetMetric() {
	//	//var tmp frrProto.OSPFInterface
	//	tmp.Area = value.GetLabel()[0].GetValue()
	//  tmp.Name = value.GetLabel()[1].GetValue()
	//	tmp.Id = value.GetLabel()[2].GetValue()
	//	tmp.Ip = value.GetLabel()[3].GetValue()
	//	neighbors = append(neighbors, &neighbor)
	//}

	return result
}

func getLsas(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFlsa {
	var result []*frrProto.OSPFlsa

	//fmt.Println(metrics["frr_ospf_neighbor_state"].GetMetric())
	//for _, value := range metrics["frr_ospf_neighbor_state"].GetMetric() {
	//	//var neighbor frrProto.OSPFInterface

	//	neighbor.Area = value.GetLabel()[0].GetValue()
	//	neighbor.Interface = value.GetLabel()[1].GetValue()
	//	neighbor.Id = value.GetLabel()[2].GetValue()
	//	neighbor.Ip = value.GetLabel()[3].GetValue()
	//	neighbors = append(neighbors, &neighbor)
	//}

	return result
}

func getRoutes(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFRoute {
	var result []*frrProto.OSPFRoute

	//fmt.Println(metrics["frr_ospf_neighbor_state"].GetMetric())
	//for _, value := range metrics["frr_ospf_neighbor_state"].GetMetric() {
	//	//var neighbor frrProto.OSPFInterface

	//	neighbor.Area = value.GetLabel()[0].GetValue()
	//	neighbor.Interface = value.GetLabel()[1].GetValue()
	//	neighbor.Id = value.GetLabel()[2].GetValue()
	//	neighbor.Ip = value.GetLabel()[3].GetValue()
	//	neighbors = append(neighbors, &neighbor)
	//}

	return result
}

func (f *Fetcher) CollectSystemMetrics() (*frrProto.SystemMetrics, error) {
	metrics := &frrProto.SystemMetrics{}

	if cpu, err := getCPUUsage(); err == nil {
		metrics.CpuUsage = cpu
	}

	if mem, err := getMemoryUsage(); err == nil {
		metrics.MemoryUsage = mem
	}

	// if stats, err := getInterfaceStats(); err == nil {
	// 	metrics.NetworkStats = stats
	// }

	return metrics, nil
}

// Helper functions
func getCPUUsage() (float64, error) {
	// This is right now only for linux
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1}'")
		out, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		var usage float64
		_, err = fmt.Sscanf(string(out), "%f", &usage)
		return usage, err
	}
	return 0, nil
}

func getMemoryUsage() (float64, error) {
	// This is right now only for linux
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'")
		out, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		var usage float64
		_, err = fmt.Sscanf(string(out), "%f", &usage)
		return usage, err
	}
	return 0, nil
}

// func getInterfaceStats() ([]frrProto.InterfaceStats, error) {
// 	// Maybe TODO
// 	return nil, nil
// }

// Functions for testing maybe remove later
func (f *Fetcher) GetMetricURLForTesting() string {
	return f.metricsURL
}

func (f *Fetcher) GetClientForTesting() *http.Client {
	return f.client
}
