**Disclaimer**

The current version of the frr-mad-analyzer is only compatible with FRR 8.5.4.
Support for newer FRR releases is planned for the near future.

---

# FRR-MAD

FRR-MAD (Free Range Routing – Monitoring and Anomaly Detection) consists of two tools, frr-mad-analyzer and frr-mad-tui. The analyzer component analyzes the static Free Range Routing configuration and compares it against the runtime lsdb. It detects wrongly advertised routes and reports them as such. The results are exposed via a Prometheus Node Exporter layer inherent to the analyzer component.

The frr-mad-tui component is a useful text user interface to give live results from the analyzer. Anomalies are presented at the dashboard and should contain valuable information. Apart from that, the tool also provides much useful information pertaining to OSPF.

| ![How the feature works](https://github.com/FRR-MAD/frr-mad-assets/blob/main/demo/demo.gif) |
|:----------------------------------------------:|
|     For more demos see [DEMO.md](DEMO.md)      |

## Introduction

This Project is split into two parts:
- **frr-mad-analyzer**: The analysis system that consists of aggregator, analyzer, exporting and socket. It spawns a Unix socket, which is accessed by frr-mad-tui, to fetch all available data. The exporter collects routing data and anomalies, exports them via the well-defined Prometheus Node Exporter API.
- **frr-mad-tui**: frr-mad-tui is the frontend of this project. It's optional but highly practical. It enables swifter sanity checks of OSPF, by providing the most useful information neatly displayed. A dynamic filter function provides additional Quality of Life improvements to the experience. Regardless of the expertise of the user, it's a useful tool to get a quick run-down of an OSPF system.

## Usage

The backend application features a handy help output. Executing the application without any arguments provides a list of available commands.
```sh
r101:/app# frr-mad-analyzer
A CLI tool for managing the FRR-MAD application.

Usage:
  frr-mad-analyzer [command]

Available Commands:
  help        Help about any command
  restart     Restart the FRR-MAD application
  start       Start the FRR-MAD application
  stop        Stop the FRR-MAD application
  version     show version number and exit

Flags:
  -h, --help   help for analyzer_frr

Use "frr-mad-analyzer [command] --help" for more information about a command.
```

### Starting the Application

To start either application a configuration file needs to be provided. Below are two options to do so. Further information for advanced settings is available in the build instructions.
- **FRR_MAD_CONFFILE Env**: Export **FRR_MAD_CONFFILE** with the absolute path to the configuration file. frr-mad-tui will use the file specified in the environment variable, otherwise it will default to /etc/frr-mad/main.yaml. The frr-mad-analyzer service works regardless of what the environment variable is.
  - This could be set with /etc/environment or /etc/profile, pick your poison.
- **--configFile Option**: When starting the daemon a custom configuration file can be provided. The path can be absolute or relative. Otherwise it will default to /etc/frr-mad/main.yaml as well.

To run the tui the environment variable is **mandatory**. If neither is provided the applications will both default to **/etc/frr-mad/main.yaml**.

#### Daemon
```sh
/path/to/frr-mad-analyzer  start --configFile /path/to/configuration
```

#### Frontend
```sh
export FRR_MAD_CONFFILE=/path/to/configuration
/path/to/frr-mad-tui
```

## Build

It's recommended to have a dedicated build host for frr-mad. The applications should be built statically, to remove any dependency issues. To build it, clone the repo and execute make. Provided make is installed. Otherwise follow the build instructions down below.

```sh
mkdir -p /tmp/frr-mad-binaries/
cd src/backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o /tmp/frr-mad-binaries/frr-mad-analyzer ./cmd/frr-analyzer
cd ../../
cd src/frontend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o /tmp/frr-mad-binaries/frr-mad-tui ./cmd/tui
cd ../../
```

### Custom Configuration Path
The default configuration path can be overridden during the build process with build flags. Provide the build flag `-X configs.ConfigLocation=/path/to/configuration.yaml` to the `go build` command.

```sh
mkdir -p /tmp/frr-mad-binaries/
cd src/backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -X configs.ConfigLocation=/path/to/configuration.yaml' -o /tmp/frr-mad-binaries/frr-mad-analyzer ./cmd/frr-analyzer
cd ../../
cd src/frontend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -X configs.ConfigLocation=/path/to/configuration.yaml' -o /tmp/frr-mad-binaries/frr-mad-tui ./cmd/tui
cd ../../
```

## Configuration File Example
```sh
mkdir -p /etc/frr-mad
cat <<EOF>/etc/frr-mad/main.yaml
default:
  tempfiles: /tmp/frr-mad
  exportpath: /tmp/frr-mad/exports
  logpath: /var/log/frr-mad
  # set debugLevel to receive different levels of logging
  # debug > info >  warn > error > none
  # Debug provides the most verbose output but it's highly resource intensive. The default is error.
  #debuglevel: error

frrmadtui:
  pages:
    ospf:
      enabled: true
    rib:
      enabled: true
    shell:
      enabled: true

socket:
  unixsocketlocation: /var/run/frr-mad
  unixsocketname: analyzer.sock
  sockettype: unix

aggregator:
  frrconfigpath: /etc/frr/frr.conf
  pollinterval: 5
  socketpath: /var/run/frr

exporter:
  # default: Port: 9091
  OSPFRouterData: true
  OSPFNetworkData: true
  OSPFSummaryData: false
  OSPFAsbrSummaryData: false
  OSPFExternalData: false
  OSPFNssaExternalData: false
  OSPFDatabase: false
  OSPFNeighbors: false
  InterfaceList: false
  RouteList: false

EOF
```

---

Refere to [Docs](docs) for more project related information
