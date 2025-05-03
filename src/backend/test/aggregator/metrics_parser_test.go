package aggregator_test

import (
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseOSPFRouterLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFRouterData
		wantErr  bool
	}{
		{
			name: "basic router LSA",
			input: `{
				"routerId": "1.1.1.1",
				"Router Link States": {
					"0.0.0.0": {
						"1.1.1.1": {
							"lsaAge": "3600",
							"lsaType": "Router",
							"linkStateId": "1.1.1.1",
							"advertisingRouter": "1.1.1.1"
						}
					}
				}
			}`,
			expected: &frrProto.OSPFRouterData{
				RouterId: "1.1.1.1",
				RouterStates: map[string]*frrProto.OSPFRouterArea{
					"0.0.0.0": {
						LsaEntries: map[string]*frrProto.OSPFRouterLSA{
							"1.1.1.1": {
								LsaAge:            3600,
								LsaType:           "Router",
								LinkStateId:       "1.1.1.1",
								AdvertisingRouter: "1.1.1.1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty input",
			input:    `{}`,
			expected: &frrProto.OSPFRouterData{},
			wantErr:  false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFRouterLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.RouterStates != nil {
				for areaID, expectedArea := range tt.expected.RouterStates {
					resultArea, ok := result.RouterStates[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					for lsaID, expectedLSA := range expectedArea.LsaEntries {
						resultLSA, ok := resultArea.LsaEntries[lsaID]
						assert.True(t, ok, "missing LSA %s", lsaID)

						assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
						assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
						assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
						assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
					}
				}
			}
		})
	}
}

func TestParseOSPFNetworkLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFNetworkData
		wantErr  bool
	}{
		{
			name: "basic network LSA",
			input: `{
				"routerId": "1.1.1.1",
				"Net Link States": {
					"0.0.0.0": {
						"2.2.2.2": {
							"lsaAge": "1800",
							"lsaType": "Network",
							"linkStateId": "2.2.2.2",
							"advertisingRouter": "1.1.1.1"
						}
					}
				}
			}`,
			expected: &frrProto.OSPFNetworkData{
				RouterId: "1.1.1.1",
				NetStates: map[string]*frrProto.NetAreaState{
					"0.0.0.0": {
						LsaEntries: map[string]*frrProto.NetworkLSA{
							"2.2.2.2": {
								LsaAge:            1800,
								LsaType:           "Network",
								LinkStateId:       "2.2.2.2",
								AdvertisingRouter: "1.1.1.1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFNetworkLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.NetStates != nil {
				for areaID, expectedArea := range tt.expected.NetStates {
					resultArea, ok := result.NetStates[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					for lsaID, expectedLSA := range expectedArea.LsaEntries {
						resultLSA, ok := resultArea.LsaEntries[lsaID]
						assert.True(t, ok, "missing LSA %s", lsaID)

						assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
						assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
						assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
						assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
					}
				}
			}
		})
	}
}

func TestParseOSPFSummaryLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFSummaryData
		wantErr  bool
	}{
		{
			name: "basic summary LSA",
			input: `{
				"routerId": "1.1.1.1",
				"Summary Link States": {
					"0.0.0.0": {
						"3.3.3.3": {
							"lsaAge": "1200",
							"lsaType": "Summary",
							"linkStateId": "3.3.3.3",
							"advertisingRouter": "1.1.1.1"
						}
					}
				}
			}`,
			expected: &frrProto.OSPFSummaryData{
				RouterId: "1.1.1.1",
				SummaryStates: map[string]*frrProto.SummaryAreaState{
					"0.0.0.0": {
						LsaEntries: map[string]*frrProto.SummaryLSA{
							"3.3.3.3": {
								LsaAge:            1200,
								LsaType:           "Summary",
								LinkStateId:       "3.3.3.3",
								AdvertisingRouter: "1.1.1.1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFSummaryLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.SummaryStates != nil {
				for areaID, expectedArea := range tt.expected.SummaryStates {
					resultArea, ok := result.SummaryStates[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					for lsaID, expectedLSA := range expectedArea.LsaEntries {
						resultLSA, ok := resultArea.LsaEntries[lsaID]
						assert.True(t, ok, "missing LSA %s", lsaID)

						assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
						assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
						assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
						assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
					}
				}
			}
		})
	}
}

func TestParseOSPFAsbrSummaryLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFAsbrSummaryData
		wantErr  bool
	}{
		{
			name: "basic ASBR summary LSA",
			input: `{
				"routerId": "1.1.1.1",
				"ASBR-Summary Link States": {
					"0.0.0.0": {
						"4.4.4.4": {
							"lsaAge": "1500",
							"lsaType": "ASBR-Summary",
							"linkStateId": "4.4.4.4",
							"advertisingRouter": "1.1.1.1"
						}
					}
				}
			}`,
			expected: &frrProto.OSPFAsbrSummaryData{
				RouterId: "1.1.1.1",
				AsbrSummaryStates: map[string]*frrProto.SummaryAreaState{
					"0.0.0.0": {
						LsaEntries: map[string]*frrProto.SummaryLSA{
							"4.4.4.4": {
								LsaAge:            1500,
								LsaType:           "ASBR-Summary",
								LinkStateId:       "4.4.4.4",
								AdvertisingRouter: "1.1.1.1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFAsbrSummaryLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.AsbrSummaryStates != nil {
				for areaID, expectedArea := range tt.expected.AsbrSummaryStates {
					resultArea, ok := result.AsbrSummaryStates[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					for lsaID, expectedLSA := range expectedArea.LsaEntries {
						resultLSA, ok := resultArea.LsaEntries[lsaID]
						assert.True(t, ok, "missing LSA %s", lsaID)

						assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
						assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
						assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
						assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
					}
				}
			}
		})
	}
}

func TestParseOSPFExternalLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFExternalData
		wantErr  bool
	}{
		{
			name: "basic external LSA",
			input: `{
				"routerId": "1.1.1.1",
				"AS External Link States": {
					"5.5.5.5": {
						"lsaAge": "2000",
						"lsaType": "AS-External",
						"linkStateId": "5.5.5.5",
						"advertisingRouter": "1.1.1.1"
					}
				}
			}`,
			expected: &frrProto.OSPFExternalData{
				RouterId: "1.1.1.1",
				AsExternalLinkStates: map[string]*frrProto.ExternalLSA{
					"5.5.5.5": {
						LsaAge:            2000,
						LsaType:           "AS-External",
						LinkStateId:       "5.5.5.5",
						AdvertisingRouter: "1.1.1.1",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFExternalLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.AsExternalLinkStates != nil {
				for lsaID, expectedLSA := range tt.expected.AsExternalLinkStates {
					resultLSA, ok := result.AsExternalLinkStates[lsaID]
					assert.True(t, ok, "missing LSA %s", lsaID)

					assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
					assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
					assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
					assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
				}
			}
		})
	}
}

func TestParseOSPFNssaExternalLSA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFNssaExternalData
		wantErr  bool
	}{
		{
			name: "basic NSSA external LSA",
			input: `{
				"routerId": "1.1.1.1",
				"NSSA-external Link States": {
					"0.0.0.1": {
						"6.6.6.6": {
							"lsaAge": "2500",
							"lsaType": "NSSA-External",
							"linkStateId": "6.6.6.6",
							"advertisingRouter": "1.1.1.1"
						}
					}
				}
			}`,
			expected: &frrProto.OSPFNssaExternalData{
				RouterId: "1.1.1.1",
				NssaExternalLinkStates: map[string]*frrProto.NssaExternalArea{
					"0.0.0.1": {
						Data: map[string]*frrProto.NssaExternalLSA{
							"6.6.6.6": {
								LsaAge:            2500,
								LsaType:           "NSSA-External",
								LinkStateId:       "6.6.6.6",
								AdvertisingRouter: "1.1.1.1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFNssaExternalLSA([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.NssaExternalLinkStates != nil {
				for areaID, expectedArea := range tt.expected.NssaExternalLinkStates {
					resultArea, ok := result.NssaExternalLinkStates[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					for lsaID, expectedLSA := range expectedArea.Data {
						resultLSA, ok := resultArea.Data[lsaID]
						assert.True(t, ok, "missing LSA %s", lsaID)

						assert.Equal(t, expectedLSA.LsaAge, resultLSA.LsaAge)
						assert.Equal(t, expectedLSA.LsaType, resultLSA.LsaType)
						assert.Equal(t, expectedLSA.LinkStateId, resultLSA.LinkStateId)
						assert.Equal(t, expectedLSA.AdvertisingRouter, resultLSA.AdvertisingRouter)
					}
				}
			}
		})
	}
}

func TestParseFullOSPFDatabase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFDatabase
		wantErr  bool
	}{
		{
			name: "basic OSPF database",
			input: `{
				"routerId": "1.1.1.1",
				"areas": {
					"0.0.0.0": {
						"routerLinkStates": [
							{
								"lsId": "1.1.1.1",
								"advertisedRouter": "1.1.1.1",
								"lsaAge": "3600"
							}
						]
					}
				}
			}`,
			expected: &frrProto.OSPFDatabase{
				RouterId: "1.1.1.1",
				Areas: map[string]*frrProto.OSPFDatabaseArea{
					"0.0.0.0": {
						RouterLinkStates: []*frrProto.RouterDataLSA{
							{
								Base: &frrProto.BaseLSA{
									LsId:             "1.1.1.1",
									AdvertisedRouter: "1.1.1.1",
									LsaAge:           3600,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseFullOSPFDatabase([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.Areas != nil {
				for areaID, expectedArea := range tt.expected.Areas {
					resultArea, ok := result.Areas[areaID]
					assert.True(t, ok, "missing area %s", areaID)

					assert.Equal(t, len(expectedArea.RouterLinkStates), len(resultArea.RouterLinkStates))
					for i, expectedLSA := range expectedArea.RouterLinkStates {
						assert.Equal(t, expectedLSA.Base.LsId, resultArea.RouterLinkStates[i].Base.LsId)
						assert.Equal(t, expectedLSA.Base.AdvertisedRouter, resultArea.RouterLinkStates[i].Base.AdvertisedRouter)
						assert.Equal(t, expectedLSA.Base.LsaAge, resultArea.RouterLinkStates[i].Base.LsaAge)
					}
				}
			}
		})
	}
}

func TestParseOSPFDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFDuplicates
		wantErr  bool
	}{
		{
			name: "basic OSPF duplicates",
			input: `{
				"routerId": "1.1.1.1",
				"asExternalLinkStates": [
					{
						"lsaAge": "3000",
						"lsaType": "AS-External",
						"linkStateId": "7.7.7.7",
						"advertisingRouter": "1.1.1.1"
					}
				]
			}`,
			expected: &frrProto.OSPFDuplicates{
				RouterId: "1.1.1.1",
				AsExternalLinkStates: []*frrProto.ASExternalLinkState{
					{
						LsaAge:            3000,
						LsaType:           "AS-External",
						LinkStateId:       "7.7.7.7",
						AdvertisingRouter: "1.1.1.1",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFDuplicates([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.RouterId, result.RouterId)

			if tt.expected.AsExternalLinkStates != nil {
				assert.Equal(t, len(tt.expected.AsExternalLinkStates), len(result.AsExternalLinkStates))
				for i, expectedLSA := range tt.expected.AsExternalLinkStates {
					assert.Equal(t, expectedLSA.LsaAge, result.AsExternalLinkStates[i].LsaAge)
					assert.Equal(t, expectedLSA.LsaType, result.AsExternalLinkStates[i].LsaType)
					assert.Equal(t, expectedLSA.LinkStateId, result.AsExternalLinkStates[i].LinkStateId)
					assert.Equal(t, expectedLSA.AdvertisingRouter, result.AsExternalLinkStates[i].AdvertisingRouter)
				}
			}
		})
	}
}

func TestParseOSPFNeighbors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.OSPFNeighbors
		wantErr  bool
	}{
		{
			name: "basic OSPF neighbors",
			input: `{
				"neighbors": {
					"eth0": [
						{
							"priority": "1",
							"state": "Full",
							"address": "10.0.0.2"
						}
					]
				}
			}`,
			expected: &frrProto.OSPFNeighbors{
				Neighbors: map[string]*frrProto.NeighborList{
					"eth0": {
						Neighbors: []*frrProto.Neighbor{
							{
								Priority: 1,
								State:    "Full",
								Address:  "10.0.0.2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseOSPFNeighbors([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.expected.Neighbors != nil {
				for iface, expectedNeighbors := range tt.expected.Neighbors {
					resultNeighbors, ok := result.Neighbors[iface]
					assert.True(t, ok, "missing interface %s", iface)

					assert.Equal(t, len(expectedNeighbors.Neighbors), len(resultNeighbors.Neighbors))
					for i, expectedNeighbor := range expectedNeighbors.Neighbors {
						assert.Equal(t, expectedNeighbor.Priority, resultNeighbors.Neighbors[i].Priority)
						assert.Equal(t, expectedNeighbor.State, resultNeighbors.Neighbors[i].State)
						assert.Equal(t, expectedNeighbor.Address, resultNeighbors.Neighbors[i].Address)
					}
				}
			}
		})
	}
}

func TestParseInterfaceStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.InterfaceList
		wantErr  bool
	}{
		{
			name: "basic interface status",
			input: `{
				"eth0": {
					"administrativeStatus": "up",
					"operationalStatus": "up",
					"index": 1
				}
			}`,
			expected: &frrProto.InterfaceList{
				Interfaces: map[string]*frrProto.SingleInterface{
					"eth0": {
						AdministrativeStatus: "up",
						OperationalStatus:    "up",
						Index:                1,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseInterfaceStatus([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.expected.Interfaces != nil {
				for ifaceName, expectedIface := range tt.expected.Interfaces {
					resultIface, ok := result.Interfaces[ifaceName]
					assert.True(t, ok, "missing interface %s", ifaceName)

					assert.Equal(t, expectedIface.AdministrativeStatus, resultIface.AdministrativeStatus)
					assert.Equal(t, expectedIface.OperationalStatus, resultIface.OperationalStatus)
					assert.Equal(t, expectedIface.Index, resultIface.Index)
				}
			}
		})
	}
}

func TestParseRib(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *frrProto.RoutingInformationBase
		wantErr  bool
	}{
		{
			name: "basic RIB",
			input: `{
				"10.0.0.0/24": [
					{
						"prefix": "10.0.0.0",
						"prefixLen": 24,
						"protocol": "ospf"
					}
				]
			}`,
			expected: &frrProto.RoutingInformationBase{
				Routes: map[string]*frrProto.RouteEntry{
					"10.0.0.0/24": {
						Routes: []*frrProto.Route{
							{
								Prefix:    "10.0.0.0",
								PrefixLen: 24,
								Protocol:  "ospf",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregator.ParseRib([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.expected.Routes != nil {
				for prefix, expectedEntry := range tt.expected.Routes {
					resultEntry, ok := result.Routes[prefix]
					assert.True(t, ok, "missing prefix %s", prefix)

					assert.Equal(t, len(expectedEntry.Routes), len(resultEntry.Routes))
					for i, expectedRoute := range expectedEntry.Routes {
						assert.Equal(t, expectedRoute.Prefix, resultEntry.Routes[i].Prefix)
						assert.Equal(t, expectedRoute.PrefixLen, resultEntry.Routes[i].PrefixLen)
						assert.Equal(t, expectedRoute.Protocol, resultEntry.Routes[i].Protocol)
					}
				}
			}
		})
	}
}
