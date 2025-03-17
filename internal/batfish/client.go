package batfish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	batfishHost = "http://localhost:9997"
)

type BatfishClient struct {
	client *http.Client
}

func NewBatfishClient() *BatfishClient {
	return &BatfishClient{
		client: &http.Client{},
	}
}

// UploadSnapshot uploads a network snapshot to Batfish
func (c *BatfishClient) UploadSnapshot(snapshotPath, snapshotName string) error {
	// Create a zip of the snapshot directory
	zipPath := filepath.Join(snapshotPath, "snapshot.zip")
	if err := zipDirectory(snapshotPath, zipPath); err != nil {
		return fmt.Errorf("failed to zip snapshot: %v", err)
	}

	// Open the zip file
	file, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open snapshot zip: %v", err)
	}
	defer file.Close()

	// Upload the snapshot
	url := fmt.Sprintf("%s/v2/snapshots", batfishHost)
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/zip")
	req.Header.Set("Snapshot-Name", snapshotName)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload snapshot: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload snapshot: %s", string(body))
	}

	return nil
}

// RunAnalysis runs a Batfish analysis and returns the results
func (c *BatfishClient) RunAnalysis(snapshotName, question string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v2/questions", batfishHost)
	payload := map[string]string{
		"snapshotName": snapshotName,
		"question":     question,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to run analysis: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to run analysis: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}

// Helper function to zip a directory
func zipDirectory(source, target string) error {
	// Implement zip logic or use a library like `archive/zip`
	return nil
}
