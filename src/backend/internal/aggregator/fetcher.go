package aggregator

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
