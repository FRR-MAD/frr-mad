package analyzer

import (
	"time"
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
