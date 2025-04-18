package analyzer

import (
	"time"
	// "github.com/ba2025-ysmprc/frr-mad/src/logger"
)

/*

 */

func InitAnalyzer(config map[string]string) *Analyzer {

	return newAnalyzer()
}

func StartAnalyzer(analyzer *Analyzer, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			//_, err := analyzer.Analyze()
			//if err != nil {
			//	log.Printf("Wollection err: %v", err)
			//	continue
			//}
		}
	}()

}

//func analyzerLogger() {
//	// Initialize log manager
//	logManager, err := logger.NewLogManager("logs/backend_app.log", logger.INFO)
//	if err != nil {
//		log.Fatalf("Failed to initialize logging: %v", err)
//	}
//
//	// Register service loggers
//	err = logManager.RegisterServiceLogger("ospf", "logs/ospf_service.log", logger.DEBUG)
//	if err != nil {
//		log.Fatalf("Failed to initialize OSPF logging: %v", err)
//	}
//
//	err = logManager.RegisterServiceLogger("bgp", "logs/bgp_service.log", logger.DEBUG)
//	if err != nil {
//		log.Fatalf("Failed to initialize BGP logging: %v", err)
//	}
//
//	// Get application logger
//	appLogger := logManager.GetAppLogger()
//	appLogger.Info("Backend application started")
//
//	// Get service logger
//	ospfLogger, ok := logManager.GetServiceLogger("ospf")
//	if !ok {
//		appLogger.Error("OSPF logger not found")
//	} else {
//		// Use standard Logger methods
//		ospfLogger.Debug("OSPF service initializing")
//
//		// Use ServiceLogger-specific methods
//		ospfLogger.LogServiceEvent("startup", "OSPF process initialized with default config")
//	}
//
//	// Continue with application...
//}
//
