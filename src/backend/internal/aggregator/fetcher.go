package aggregator

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"

	frrSocket "github.com/frr-mad/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

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

func fetchGeneralOSPFInformation(executor *frrSocket.FRRCommandExecutor) (*frrProto.GeneralOspfInformation, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf json")
	if err != nil {
		return nil, err
	}

	return ParseGeneralOspfInformation(output)
}

func fetchOSPFRouterData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFRouterData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data router self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFRouterLSA(output)
}

func fetchOSPFRouterDataAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFRouterData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data router json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFRouterLSAAll(output)
}

func fetchOSPFNetworkData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNetworkData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data network self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNetworkLSA(output)
}

func fetchOSPFNetworkDataAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNetworkData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data network json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFNetworkLSAAll(output)
}

func fetchOSPFSummaryData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFSummaryLSA(output)
}

func fetchOSPFSummaryDataAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data summary json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFSummaryLSAAll(output)
}

func fetchOSPFAsbrSummaryData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFAsbrSummaryData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data asbr-summary self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFAsbrSummaryLSA(output)
}

func fetchOSPFExternalData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFExternalLSA(output)
}

func fetchOSPFNssaExternalData(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNssaExternalData, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data nssa-external self json")
	if err != nil {
		return nil, err
	}

	return ParseOSPFNssaExternalLSA(output)
}

func fetchFullOSPFDatabase(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFDatabase, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data json")
	if err != nil {
		return nil, err
	}
	return ParseFullOSPFDatabase(output)
}

func fetchOSPFExternalAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFExternalAll, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data external json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFExternalAll(output)
}

func fetchOSPFNssaExternalAll(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNssaExternalAll, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf data nssa-external json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFNssaExternalAll(output)
}

func fetchOSPFNeighbors(executor *frrSocket.FRRCommandExecutor) (*frrProto.OSPFNeighbors, error) {
	output, err := executor.ExecOSPFCmd("show ip ospf neighbor json")
	if err != nil {
		return nil, err
	}
	return ParseOSPFNeighbors(output)
}

func fetchInterfaceStatus(executor *frrSocket.FRRCommandExecutor) (*frrProto.InterfaceList, error) {
	output, err := executor.ExecZebraCmd("show interface json")
	if err != nil {
		return nil, err
	}
	return ParseInterfaceStatus(output)
}

func fetchRib(executor *frrSocket.FRRCommandExecutor) (*frrProto.RoutingInformationBase, error) {
	output, err := executor.ExecZebraCmd("show ip route json")
	if err != nil {
		return nil, err
	}
	return ParseRib(output)
}

func fetchRibFibSummary(executor *frrSocket.FRRCommandExecutor) (*frrProto.RibFibSummaryRoutes, error) {
	output, err := executor.ExecZebraCmd("show ip route summary json")
	if err != nil {
		return nil, err
	}
	return ParseRibFibSummary(output)
}

func collectSystemMetrics() (*frrProto.SystemMetrics, error) {
	metrics := &frrProto.SystemMetrics{}

	cores, err := getCPUAmount()
	if err == nil {
		metrics.CpuAmount = cores
	}

	if cpu, err := getCPUUsagePercent(); err == nil {
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

func getCPUUsagePercent() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) == 0 {
		return 0, fmt.Errorf("no CPU usage data returned")
	}

	// First val is avg
	usage := percentages[0]

	if usage < 0 {
		usage = 0
	} else if usage > 100 {
		usage = 100
	}

	return usage / 100, nil
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
