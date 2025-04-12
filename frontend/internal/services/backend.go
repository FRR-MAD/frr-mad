package backend

func GetOSPFAnomalies() [][]string {
	// Fetch OSPF Anomalies via protobuf

	// parse received protobuf data

	// parsed protobuf message should look something like this:
	anomalyRows := [][]string{
		{"10.0.12.0/23", "unadvertised route", "OSPF Monitoring Tab 5", "Start"},
		{"10.0.15.0/14", "wrongly advertised", "OSPF Monitoring Tab 3", "Start"},
		{"10.0.199.0/23", "overadvertised route", "OSPF Monitoring Tab 2", "Start"},
	}

	return anomalyRows

}

func GetOSPFMetrics() [][]string {
	// Fetch all metrics (maybe fetch periodically everything and with the Getter function only provide requested data

	// this getter provides the OSPF metrics for the dashboard if no anomaly is detected

	allGoodRows := [][]string{
		{"10.0.15.0/24", "Type 5 - AS External LSA", "correct"},
		{"10.0.16.0/24", "Type 5 - AS External LSA", "correct"},
		{"10.0.17.0/24", "Type 5 - AS External LSA", "correct"},
		{"10.0.18.0/24", "Type 5 - AS External LSA", "correct"},
	}

	return allGoodRows
}
