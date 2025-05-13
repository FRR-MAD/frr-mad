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

func FetchGeneralOSPFInformation(executor *frrSocket.FRRCommandExecutor) (*frrProto.GeneralOspfInformation, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf json")
	if err != nil {
		return nil, err
	}

	return ParseGeneralOspfInformation(output)
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

func FetchOSPFExternalAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFExternalAll, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf database external json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFExternalAll(output)
}

func FetchOSPFNssaExternalAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNssaExternalAll, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf database nssa-external json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFNssaExternalAll(output)
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

func FetchRib(executor *frrSocket.FRRCommandExecutor) (*frrProto.RoutingInformationBase, error) {
	output, err := executor.ExecZebraCmd("show ip route json")
	if err != nil {
		return nil, err
	}
	return ParseRib(output)
}

func FetchRibFibSummary(executor *frrSocket.FRRCommandExecutor) (*frrProto.RibFibSummaryRoutes, error) {
	output, err := executor.ExecZebraCmd("show ip route summary json")
	if err != nil {
		return nil, err
	}
	return ParseRibFibSummary(output)
}

func (f *Fetcher) CollectSystemMetrics() (*frrProto.SystemMetrics, error) {
	metrics := &frrProto.SystemMetrics{}

	cores, err := getCPUAmount()
	if err == nil {
		metrics.CpuAmount = cores
	}

	if cpu, err := getCPUUsagePercent(1*time.Second, int(cores)); err == nil {
		metrics.CpuUsage = cpu
	}

	if mem, err := getMemoryUsage(); err == nil {
		metrics.MemoryUsage = mem
	}

	return metrics, nil
}

// Helper functions
func getCPUAmount() (int64, error) {
	cores := runtime.NumCPU()
	return int64(cores), nil
}

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
	var values []float64
	for _, s := range fields[1:] {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, 0, err
		}
		values = append(values, v)
	}

	idle = values[3]
	for _, v := range values {
		total += v
	}
	return total, idle, nil
}

func getCPUUsagePercent(interval time.Duration, cores int) (float64, error) {
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
	// Normalize by number of cores
	normalizedBusy := busy / float64(cores)

	fmt.Println(busy)
	fmt.Println(normalizedBusy)

	return normalizedBusy, nil
}

func getMemoryUsage() (float64, error) {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/meminfo")
		if err != nil {
			return 0, err
		}
		lines := strings.Split(string(data), "\n")
		var total, free, buffers, cached uint64
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				total = parseMemLine(line)
			} else if strings.HasPrefix(line, "MemFree:") {
				free = parseMemLine(line)
			} else if strings.HasPrefix(line, "Buffers:") {
				buffers = parseMemLine(line)
			} else if strings.HasPrefix(line, "Cached:") {
				cached = parseMemLine(line)
			}
		}
		if total == 0 {
			return 0, fmt.Errorf("could not read memory stats")
		}
		used := total - free - buffers - cached
		return float64(used) / float64(total) * 100, nil
	}
	return 0, nil
}

func parseMemLine(line string) uint64 {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0
	}
	val, _ := strconv.ParseUint(parts[1], 10, 64)
	return val
}
