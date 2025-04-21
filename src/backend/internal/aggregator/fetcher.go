package aggregator

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
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

func fetchStaticFRRConfig() (*frrProto.StaticFRRConfiguration, error) {
	cmd := exec.Command("vtysh", "-c", "show running-config")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can not open file: %w", err)
	}

	tmp, err := os.Create("/tmp/frr-config.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmp.Write(output); err != nil {
		err := tmp.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to write config to temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	parsedStaticFRRConfig, err := ParseStaticFRRConfig(tmp.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to parse static FRR config: %w", err)
	}

	return parsedStaticFRRConfig, nil
}

func FetchOSPFRouterData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFRouterData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data router self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFRouterLSA(output)
}

func FetchOSPFNetworkData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNetworkData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data network self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNetworkLSA(output)
}

func FetchOSPFSummaryData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFSummaryLSA(output)
}

func FetchOSPFAsbrSummaryData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFAsbrSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data asbr-summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFAsbrSummaryLSA(output)
}

func FetchOSPFExternalData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFExternalLSA(output)
}

func FetchOSPFNssaExternalData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNssaExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data nssa-external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNssaExternalLSA(output)
}

func FetchFullOSPFDatabase(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFDatabase, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf database json")
	if err != nil {
		return nil, err
	}
	return ParseFullOSPFDatabase(output)
}

func FetchOSPFDuplicateCandidates(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFDuplicates, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf database external json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFDuplicates(output)
}

func FetchOSPFNeighbors(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNeighbors, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf neighbor json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFNeighbors(output)
}

func FetchInterfaceStatus(executor *frrSocket.FRRCommandExecutor) (*frrProto.InterfaceList, error) {
	output, err := executor.ExecZebraCmd("show interface json")
	if err != nil {
		return nil, err
	}
	return ParseInterfaceStatus(output)
}

func FetchExpectedRoutes(executor *frrSocket.FRRCommandExecutor) (*frrProto.RouteList, error) {
	output, err := executor.ExecZebraCmd("show ip route json")
	if err != nil {
		return nil, err
	}
	return ParseRouteList(output)
}

func (f *Fetcher) CollectSystemMetrics() (*frrProto.SystemMetrics, error) {
	metrics := &frrProto.SystemMetrics{}

	if cores, err := getCPUAmount(); err == nil {
		metrics.CpuAmount = cores
	}

	if cpu, err := getCPUUsagePercent(1 * time.Second); err == nil {
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
func getCPUAmount() (int64, error) {
	cores := runtime.NumCPU()
	return int64(cores), nil
}

// readCPUSample parses the first line of /proc/stat and returns
// totalJiffies and idleJiffies.
func readCPUSample() (total, idle float64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0, 0, fmt.Errorf("failed to scan /proc/stat")
	}
	fields := strings.Fields(scanner.Text())
	// fields[0]=="cpu", then user, nice, system, idle, iowait, irq, ...
	var values []float64
	for _, s := range fields[1:] {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, 0, err
		}
		values = append(values, v)
	}

	// idle is the 4th value (index 3)
	idle = values[3]
	for _, v := range values {
		total += v
	}
	return total, idle, nil
}

// getCPUUsagePercent reads two samples 'interval' apart and returns
// busy percentage in [0.0, 100.0].
func getCPUUsagePercent(interval time.Duration) (float64, error) {
	t0, id0, err := readCPUSample()
	if err != nil {
		return 0, err
	}
	time.Sleep(interval)
	t1, id1, err := readCPUSample()
	if err != nil {
		return 0, err
	}

	totalDelta := t1 - t0
	idleDelta := id1 - id0
	if totalDelta <= 0 {
		return 0, fmt.Errorf("invalid CPU delta: %.0f", totalDelta)
	}

	busy := (totalDelta - idleDelta) / totalDelta * 100.0
	// clamp to [0,100]
	if busy < 0 {
		busy = 0
	} else if busy > 100 {
		busy = 100
	}
	return busy, nil
}

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
