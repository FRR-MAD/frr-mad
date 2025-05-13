# FRR-MAD

FRR-MAD (Free Range Routing – Monitoring and Anomaly Detection) is an intuitive Terminal User Interface for monitoring OSPF states within FRRouting.
It effectively detects anomalies by comparing static file data with the Link-State Database (LSDB) and the Forwarding Information Base (FIB).

## Usage

This Project is split into two parts:
- **frr-tui**: The frontend of our application. It's not really necessary, but makes it a lot easier to check the sanity of the application. It also provides many OSPF stats, helpful to less experienced network engineers. Regardless of experience, anomalies are monitored with the frr-tui.
- **frr-analyzer**: The analysis system that consits of aggregation, analysis and information exporting. It spawns a socket, which the frr-tui unit uses to fetch all necessary data. The exporter collects routing data and exports them via the well-defined Prometheus Node Exporter uri. 

### frr-analyzer
Install the applications where you want or create your custom start scripts. To start the analyzer simply execute:
```sh
/path/to/frr-analyzer start
```

This will start a persistent service and spawn a socket. To stop it, you can execute the binary again with the stop argument. Don't kill it if possible, as it will leave a spawned socket.

By default the exporter runs on port 9091 an all ips. The port is freely adjustable.

### frr-tui
The tui can be started by running:
```
/path/to/frr-tui
```

You will enter a small terminal user interface. It serves as a small monitoring application with the most important information already present. 

## Installation

Installation is fairly easy. Clone the repo and build it. The executable is compiled with the static flag, so remove it if you have all the dependencies set on the host.

```sh
mkdir -p /tmp/frr-mad-binaries/
cd src/backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o /tmp/frr-mad-binaries/frr-mad-analyzer ./cmd/frr-analyzer
cd ../../
cd src/frontend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o /tmp/frr-mad-binaries/frr-mad-tui ./cmd/tui
cd ../../
```

Provided is a default configuration file. Creating it will start the application. 

```sh
mkdir -p /etc/frr-mad
cat <<EOF>/etc/frr-mad/frr-mad.yaml
default:
  tempfiles: /tmp/frr-mad
  logpath: /tmp/frr-mad/log
  debuglevel: none 

socket:
  unixsocketlocation: /tmp/frr-mad
  unixsocketname: analyzer.sock
  sockettype: unix

analyzer:
  foo: bar

aggregator:
  frrmetricsurl: http://localhost:9342/metrics
  frrconfigpath: /etc/frr/frr.conf
  pollinterval: 5
  socketpathbgp: /var/run/frr/bgpd.vty
  socketpathospf: /var/run/frr/ospfd.vty
  socketpathzebra: /var/run/frr/zebra.vty
  socketpath: /var/run/frr

exporter:
  Port: 9091     
  OSPFRouterData: (collector.ospf.router,"Collect OSPF router information metrics",true)
  OSPFNetworkData: (collector.ospf.network,"Collect OSPF network information metrics",true)
  OSPFSummaryData: (collector.ospf.summary,"Collect OSPF summary information metrics",true)
  OSPFAsbrSummaryData: (collector.ospf.asbr-summary,"Collect OSPF ASBR summary information metrics",true)
  OSPFExternalData: (collector.ospf.external,"Collect OSPF external route information metrics",true)
  OSPFNssaExternalData: (collector.ospf.nssa-external,"Collect OSPF NSSA external route information metrics",true)
  OSPFDatabase: (collector.ospf.database,"Collect OSPF database information metrics",true)
  OSPFDuplicates: (collector.ospf.duplicates,"Collect OSPF duplicate information metrics",true)
  OSPFNeighbors: (collector.ospf.neighbors,"Collect OSPF neighbor information metrics",true)
  InterfaceList: (collector.interface.list,"Collect interface list information metrics",true)
  RouteList: (collector.route.list,"Collect route list information metrics",true)
EOF
```

The default folders for this application are:
- config location: /etc/frr-mad/
- log location: /var/tmp/frr-mad
- tmp files: /tmp/frr-mad
- Unix socket location: /var/run/frr-mad


That's all there is to the installation and setup.

## Project Structure

```
root/
├── archive/                 # 
├── backend/                 # 
├── binaries/                # Ready to use Go binaries
├── protobufSource/          # Protofile for go-types generation
├── src/                     # Source Code 
│   ├── backend/             # 
│   │   ├── internal/        # 
│   │   │   ├── aggregator/  # Logic to fetch, process and parse data
│   │   │   ├── analyzer/    # Logic to analyze collected data
│   │   │   ├── comms/       # Unix Socket creation
│   │   │   ├── logger/      # Logic for application logging
│   ├── frontend/            # Terminal User Interface using Charmbracelet Libraries
│   └── logger/              # Project wide logger implementation using slog
└── README.md                # Project documentation
```

### Backend Aggregator Structure

```
backend
├──internal/
│  └── aggregator/          # This is your collector
│       ├── collector.go     # Main collection logic
│       ├── fetcher.go       # HTTP metrics fetching
│       ├── parser.go        # FRR config parsing
│       ├── types.go         # DTO definitions
│       └── converter.go     # Metrics to DTO conversion
```

1. **Separation by functionality**: UI, distro handling, and configuration are clearly separated
2. **Easy extension**: Adding support for a new distro only requires adding a new handler in `internal/distro/handlers/`
3. **Private vs public code**: `internal/` keeps implementation details hidden while `pkg/` exposes reusable components
4. **Single responsibility**: Each package has a clear purpose
5. **Testability**: Components are modular and can be tested independently

This is the initial design of the code environment.

### Frontend Structure

```
root/
├── src/                           # Source Code 
│   ├── frontend/                  # 
│   │   ├── cmd/                   # 
│   │   │   ├── tui/               # Entry Point (main.go)
│   │   ├── internal/              # 
│   │   │   ├── common/            # Shared types, helpers, and utilites across pages
│   │   │   ├── pages/             # Each Page has it’s own model
│   │   │   │   ├── examplePage/   #
│   │   │   │   │   ├── model/     # Bubbletea model
│   │   │   │   │   ├── update/    # update logic and message handling
│   │   │   │   │   ├── view/      # UI rendering and Backend data aggregation
│   │   │   ├── services/          # Backend service layer to call external systems
│   │   │   ├── ui/                # Shared UI styling, mainly lipgloss
```

## Development

TODO

### CI/CD

- **GitHub Actions**: For this project we will use GitHub Actions for automated nightly builds and manual deployments.