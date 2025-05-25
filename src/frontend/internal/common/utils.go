package common

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	frrProto "github.com/frr-mad/frr-tui/pkg"
)

func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func HasAnyAnomaly(a *frrProto.AnomalyDetection) bool {
	if a == nil {
		return false
	}
	return a.HasUnAdvertisedPrefixes ||
		a.HasOverAdvertisedPrefixes ||
		a.HasDuplicatePrefixes ||
		a.HasMisconfiguredPrefixes
}

func PrintBackendError(err error, functionName string) string {
	return fmt.Sprintf(
		"Error: \n%v\n\nNo data received from backend for '%s()'. Press 'r' to reload...",
		err, functionName,
	)
}

// SortTableByIPColumn sorts a 2D string slice by the first column as IP addresses.
func SortTableByIPColumn(table [][]string) {
	sort.Slice(table, func(i, j int) bool {
		ip1 := net.ParseIP(table[i][0])
		ip2 := net.ParseIP(table[j][0])

		// Fallback to lexicographic sort if parsing fails
		if ip1 == nil || ip2 == nil {
			return table[i][0] < table[j][0]
		}

		return bytesCompare(ip1.To16(), ip2.To16()) < 0
	})
}

// bytesCompare compares two byte slices representing IPs.
func bytesCompare(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return len(a) - len(b)
}

func PrettyPrintJSON(msg proto.Message) string {
	out, err := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return "Failed to marshal proto message to JSON: " + err.Error()
	}
	return string(out)
}

// ExportProto upserts an export option and always wraps the proto payload
// in a timestamped envelope.
func ExportProto(
	exportData map[string]string,
	exportOptions []ExportOption,
	key, label, filename string,
	fetch func() (proto.Message, error),
) ([]ExportOption, error) {
	msg, err := fetch()
	if err != nil {
		return exportOptions, err
	}

	// 1) marshal the proto with protojson
	dataJSON, err := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return exportOptions, err
	}

	// 2) wrap in a timestamp + RawMessage
	wrapper := struct {
		ReceivedAt string          `json:"received_at"`
		Data       json.RawMessage `json:"data"`
	}{
		ReceivedAt: time.Now().UTC().Format(time.RFC3339),
		Data:       json.RawMessage(dataJSON),
	}

	// 3) pretty-print the wrapper
	b, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return exportOptions, err
	}
	exportData[key] = string(b)

	// 4) upsert the option
	exportOptions = AddExportOption(exportOptions, ExportOption{
		Label:    label,
		MapKey:   key,
		Filename: filename,
	})
	return exportOptions, nil
}

// WriteExportToFile writes `data` into a file named `filename` under /tmp/frr-mad/exports.
// Ensures the filename ends with .json; truncates existing files.
func WriteExportToFile(data string, filename string, directory string) error {
	exportDir := directory
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	// ensure .json extension
	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	path := filepath.Join(exportDir, filename)
	// open file with truncate + create flags
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data + "\n"); err != nil {
		return err
	}
	return nil
}

func FilterRows(rows [][]string, q string) [][]string {
	if q == "" {
		return rows
	}
	var out [][]string
	for _, r := range rows {
		joined := strings.ToLower(strings.Join(r, " "))
		if strings.Contains(joined, strings.ToLower(q)) {
			out = append(out, r)
		}
	}
	return out
}
