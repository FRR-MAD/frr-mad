package analyzer

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// analyze the different ospf anomalies
// call ospf functions

type magicalState struct {
	Hostname   string `json:"hostname"`
	FrrVersion string `json:"frr_version"`
	Areas      []area `json:"areas"`
}

type area struct {
	Name    string         `json:"name"`
	LsaType string         `json:"lsa_type"`
	Links   []advertisment `json:"links"`
}

type advertisment struct {
	IpAddress    string `json:"ip_address"`
	PrefixLength uint32 `json:"prefix_length"`
	LsaType      int    `json:"lsa_type"`
	Cost         int    `json:"cost"`
}

func (c *Analyzer) AnomalyAnalysis() {
	ms := convertToMagicalState(c.metrics.StaticFrrConfiguration)

	// Print the result
	fmt.Println(ms)
	//msJSON, _ := json.MarshalIndent(ms, "", "  ")
	//fmt.Println("Magical State:")
	//fmt.Println(string(msJSON))

	//// Convert back to StaticFRRConfiguration
	//backToConfig := ConvertFromMagicalState(ms)

	//// Print the result
	//configJSON, _ := json.MarshalIndent(backToConfig, "", "  ")
	//fmt.Println("\nBack to StaticFRRConfiguration:")
	//fmt.Println(string(configJSON))

	//	tmpAreaMap := make(map[string][]advertisment)

	//for _, v := range c.metrics.StaticFrrConfiguration.Interfaces {
	//areaValue := v.GetArea()
	//if areaValue != "" {
	//for _, i := range v.IpAddress {
	//tmpAdvertisment := advertisment{
	//ip_address:    i.GetIpAddress(),
	//prefix_length: i.GetPrefixLength(),
	//lsa_type:      1,
	//cost:          10,
	//}
	//if _, exists := tmpAreaMap[areaValue]; !exists {
	//tmpAreaMap[areaValue] = make([]advertisment, 0)
	//}
	//tmpAreaMap[areaValue] = append(tmpAreaMap[areaValue], tmpAdvertisment)
	//}

	//}
	//}

	//fmt.Printf("%v\n", tmpAreaMap)
}

// ConvertToMagicalState converts a StaticFRRConfiguration to a magicalState
func convertToMagicalState(config *frrProto.StaticFRRConfiguration) *magicalState {
	result := &magicalState{
		Hostname:   config.Hostname,
		FrrVersion: config.FrrVersion,
		Areas:      []area{},
	}

	// Map to store unique areas
	areaMap := make(map[string]*area)

	// Process all interfaces
	for _, iface := range config.Interfaces {
		// Skip interfaces without an area
		if iface.Area == "" {
			continue
		}

		// Get or create area
		a, exists := areaMap[iface.Area]
		advertismentList := make([]advertisment, 0)
		if !exists {
			newArea := area{
				Name:    iface.Area,
				LsaType: "router", // Default LSA type for areas
				Links:   advertismentList,
			}
			areaMap[iface.Area] = &newArea
			a = &newArea
		}

		// Create advertisements from IP addresses
		var adv advertisment
		for _, ip := range iface.IpAddress {
			//adv := advertisement{
			//IpAddress:    ip.IpAddress,
			//PrefixLength: ip.PrefixLength,
			//LsaType:      1,
			//Cost:         10,
			//}
			adv.IpAddress = ip.IpAddress
			adv.PrefixLength = ip.PrefixLength
			adv.LsaType = 1
			adv.Cost = 10

			a.Links = append(a.Links, adv)
		}
	}

	// Convert map to slice for the final result
	for _, a := range areaMap {
		result.Areas = append(result.Areas, *a)
	}

	return result
}

// ConvertFromMagicalState converts a magicalState back to StaticFRRConfiguration//
//func ConvertFromMagicalState(ms *magicalState) *StaticFRRConfiguration {
//config := &StaticFRRConfiguration{
//Hostname:   ms.Hostname,
//FrrVersion: ms.FrrVersion,
//Interfaces: []*Interface{},
//}

//// Create a map to group IPs by area
//areaIPs := make(map[string][]*IPPrefix)

//// Process all areas and their advertisements
//for _, a := range ms.Areas {
//for _, adv := range a.Links {
//ip := &IPPrefix{
//IpAddress:    adv.IpAddress,
//PrefixLength: adv.PrefixLength,
//}
//areaIPs[a.Name] = append(areaIPs[a.Name], ip)
//}
//}

//// Create interfaces for each area
//// Note: This is a simplification. In a real scenario, you might want to
//// create more logical interface groupings based on your application's needs.
//for areaName, ips := range areaIPs {
//iface := &Interface{
//Name:      fmt.Sprintf("area_%s_interface", areaName), // Generate a name
//IpAddress: ips,
//Area:      areaName,
//Passive:   false,
//}
//config.Interfaces = append(config.Interfaces, iface)
//}

//return config
//}

//// Helper function to marshal/unmarshal
//func MarshalMagicalState(ms *magicalState) ([]byte, error) {
//return json.Marshal(ms)
//}

//func UnmarshalMagicalState(data []byte) (*magicalState, error) {
//var ms magicalState
//err := json.Unmarshal(data, &ms)
//return &ms, err
//}

func Example() {
	// Convert to magicalState
}

func (a *Analyzer) Foobar() string {
	return "mighty analyzer"
}
