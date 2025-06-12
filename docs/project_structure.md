# Project Structure

```
root/
├── protobufSource/          # Protofile for go-types generation
├── src/                     # Source Code 
│   ├── backend/             # 
│   │   ├── internal/        # 
│   │   │   ├── aggregator/  # Logic to fetch, process and parse data
│   │   │   ├── analyzer/    # Logic to analyze collected data
│   │   │   ├── socket/       # Unix Socket creation
│   │   │   ├── logger/      # Logic for application logging
│   ├── frontend/            # Terminal User Interface using Charmbracelet Libraries
│   └── logger/              # Project wide logger implementation using slog
└── README.md                # Project documentation
```

## Backend Aggregator Structure

```
backend
├──internal/
│  └── aggregator/           # This is the collector system
│       ├── collector.go     # Main collection logic
│       ├── fetcher.go       # HTTP metrics fetching
│       ├── parser.go        # FRR config parsing
│       ├── types.go         # DTO definitions
│       └── converter.go     # Metrics to DTO conversion
```
## Backend Analyzer Structure

```
backend
├──internal/
│  └── analyzer/                        # This is the analyzer system
│       ├── analyzer.go                 # Main analysis hub
│       ├── isStateLSDBParser.go        # Parses is state (lsdb)
│       ├── main.go                     # Initializes analyzer object
│       ├── ospfAnalysis.go             # Analyzes anomalies from is state and should state
│       └── shouldStateLSDBParser.go    # Parses should state from configuration
```

## Backend socket Structure

```
backend
├──internal/
│  └── socket/                        # This is the socket system 
│       ├── analysisProcessing.go     # Processing of socket calls for analyzer data
│       ├── dummyData.go              # Dummy data for testing
│       ├── frrProcessing.go          # Processing of socket calls for FRR data
│       ├── ospfProcessing.go         # Processing of socket calls for OSPF data
│       ├── processing.go             # Processing of socket calls
│       └── socket.go                 # Initializes and spawns the socket

```

## Backend exporter Structure

```
backend
├──internal/
│  └── exporter/                    # This is the exporter system 
│       ├── anomalyExporter.go      # Attaches to exporter object and exports anomaly data
│       ├── main.go                 # Initializes the exporter object
│       └── metricsExporter.go      # Attaches to exporter object and exports FRR metrics
```

# Frontend Structure

```
frontend/                  # 
├── cmd/                   # 
│   ├── tui/               # Entry Point (main.go)
├── internal/              # 
│   ├── common/            # Shared types, helpers, and utilities across pages
│   ├── pages/             # Each Page has its own model
│   │   └── examplePage/   #
│   │       ├── model/     # Bubbletea model
│   │       ├── update/    # update logic and message handling
│   │       └── view/      # UI rendering and Backend data aggregation
│   ├── services/          # Backend service layer to call external systems
│   └── ui/                # Shared UI styling, mainly lipgloss
```
