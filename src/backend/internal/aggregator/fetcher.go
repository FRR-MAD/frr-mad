package aggregator

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
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

func FetchOSPFRouterData(executor *frrSocket.FRRCommandExecutor) (*OSPFRouterData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data router self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFRouterLSA(output)
}

func FetchOSPFNetworkData(executor *frrSocket.FRRCommandExecutor) (*OSPFNetworkData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data network self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNetworkLSA(output)
}

func FetchOSPFSummaryData(executor *frrSocket.FRRCommandExecutor) (*OSPFSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFSummaryLSA(output)
}

func FetchOSPFAsbrSummaryData(executor *frrSocket.FRRCommandExecutor) (*OSPFAsbrSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data asbr-summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFAsbrSummaryLSA(output)
}

func FetchOSPFExternalData(executor *frrSocket.FRRCommandExecutor) (*OSPFExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFExternalLSA(output)
}

func FetchOSPFNssaExternalData(executor *frrSocket.FRRCommandExecutor) (*OSPFNssaExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data nssa-external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNssaExternalLSA(output)
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
	hasRouteChanges := getRouteChanges(parsedMetrics)

	metrics.Neighbors = OSPFNeighbors
	metrics.Routes = OSPFRoutes
	metrics.Interfaces = OSPFInterfaces
	metrics.Lsas = OSPFlsa
	metrics.HasRouteChanges = hasRouteChanges

	return &metrics, nil
}

func getNeighbors(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFNeighbor {
	var neighbors []*frrProto.OSPFNeighbor

	neighborMetrics, exists := metrics["frr_ospf_neighbor_state"]
	if !exists {
		return neighbors
	}

	for _, metric := range neighborMetrics.GetMetric() {
		var neighbor frrProto.OSPFNeighbor

		for _, label := range metric.GetLabel() {
			switch label.GetName() {
			case "area":
				neighbor.Area = label.GetValue()
			case "iface":
				neighbor.Interface = label.GetValue()
			case "neighbor_id":
				neighbor.Id = label.GetValue()
			case "neighbor_ip":
				neighbor.Ip = label.GetValue()
			}
		}

		if metric.Gauge != nil {
			neighbor.State = mapOSPFValueToState(int32(metric.GetGauge().GetValue()))
		}

		neighbors = append(neighbors, &neighbor)
	}

	return neighbors
}

func getLsas(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFlsa {
	var lsas []*frrProto.OSPFlsa

	lsaMetrics, exists := metrics["frr_ospf_lsa_detail"]
	if !exists {
		return lsas
	}

	for _, metric := range lsaMetrics.GetMetric() {
		var lsa frrProto.OSPFlsa

		for _, label := range metric.GetLabel() {
			switch label.GetName() {
			case "area":
				lsa.Area = label.GetValue()
			case "adv_router":
				lsa.AdvRouter = label.GetValue()
			case "lsa_id":
				lsa.LsId = label.GetValue()
			case "lsa_type":
				lsa.Type = label.GetValue()
			case "sequence":
				lsa.Sequence = label.GetValue()
			}
		}

		lsas = append(lsas, &lsa)
	}

	return lsas
}

func getRoutes(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFRoute {
	var routes []*frrProto.OSPFRoute

	routeMetrics, exists := metrics["frr_ospf_route_detail"]
	if !exists {
		return routes
	}

	for _, metric := range routeMetrics.GetMetric() {
		var route frrProto.OSPFRoute

		for _, label := range metric.GetLabel() {
			switch label.GetName() {
			case "area":
				route.Area = label.GetValue()
			case "interface":
				route.Interface = label.GetValue()
			case "next_hop":
				route.NextHop = label.GetValue()
			case "prefix":
				route.Prefix = label.GetValue()
			case "route_type":
				route.Type = label.GetValue()
			}
		}

		if metric.Gauge != nil {
			route.Cost = int32(metric.GetGauge().GetValue())
		}

		routes = append(routes, &route)
	}

	return routes
}

func getInterfaces(metrics map[string]*io_prometheus_client.MetricFamily) []*frrProto.OSPFInterface {
	var interfaces []*frrProto.OSPFInterface
	/*
			  Name string
		    Area string

		    NbrCount int32
		    NbrAdj int32
		    Passive bool
	*/
	//fmt.Println(metrics["frr_ospf_neighbor_state"].GetMetric()[1].GetLabel()[1])
	//fmt.Println(metrics["frr_ospf_neighbor_state"])
	//for _, value := range metrics["frr_ospf_neighbor_state"].GetMetric() {
	//	//var tmp frrProto.OSPFInterface
	//	tmp.Area = value.GetLabel()[0].GetValue()
	//  tmp.Name = value.GetLabel()[1].GetValue()
	//	tmp.Id = value.GetLabel()[2].GetValue()
	//	tmp.Ip = value.GetLabel()[3].GetValue()
	//	neighbors = append(neighbors, &neighbor)
	//}

	return interfaces
}

func getRouteChanges(metrics map[string]*io_prometheus_client.MetricFamily) bool {
	routeChangesMetrics, exists := metrics["frr_ospf_has_route_changes"]
	if !exists {
		return false
	}

	for _, metric := range routeChangesMetrics.GetMetric() {
		if metric.Gauge != nil && metric.GetGauge().GetValue() > 0 {
			return true
		}
	}

	return false
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

func mapOSPFValueToState(value int32) string {
	var state string
	switch int32(value) {
	case 1:
		state = "full"
	case 2:
		state = "down"
	case 3:
		state = "attempt"
	case 4:
		state = "init"
	case 5:
		state = "2way"
	case 6:
		state = "exstart"
	case 7:
		state = "exchange"
	case 8:
		state = "loading"
	default:
		state = "default"
	}

	return state
}
