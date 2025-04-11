package aggregator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"
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

func (f *Fetcher) FetchOSPF() (*OSPFMetrics, error) {
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

func parseOSPFMetrics(rawData []byte) (*OSPFMetrics, error) {
	var metrics OSPFMetrics
	if err := json.Unmarshal(rawData, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}
	return &metrics, nil
}

func (f *Fetcher) CollectSystemMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{}

	if cpu, err := getCPUUsage(); err == nil {
		metrics.CPUUsage = cpu
	}

	if mem, err := getMemoryUsage(); err == nil {
		metrics.MemoryUsage = mem
	}

	if stats, err := getInterfaceStats(); err == nil {
		metrics.NetworkStats = stats
	}

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

func getInterfaceStats() ([]InterfaceStats, error) {
	// Maybe TODO
	return nil, nil
}

// Functions for testing maybe remove later
func (f *Fetcher) GetMetricURLForTesting() string {
	return f.metricsURL
}

func (f *Fetcher) GetClientForTesting() *http.Client {
	return f.client
}
