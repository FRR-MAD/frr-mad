# frr-tui

## Installation

- foobar

## Setup Instructions for Batfish

1. Clone the repository:
    ```bash
    git clone https://github.com/ba2025-ysmprc/frr-tui.git
    cd tui/batfish
    ```
2. Start the setup script:
    ```bash
    . ./scripts/setup.sh
    ```

3. Run the python script and start the docker container for batfish:
    ```bash
    docker run --rm -it -p 9997:9997 batfish/batfish:latest
    ./run_tui.sh
    ```

## Setup Instructions for GoBGP

1. Clone the repository:
    ```bash
    git clone https://github.com/ba2025-ysmprc/frr-tui.git
    cd tui
    ```
2. Command for normal TUI
    ```bash
    go run cmd/tui/main.go
    ```

3. Command for testing TUI
    ```bash
    go run cmd/tui/main.go -test
    ```


## Project Structure

```
root/
├── cmd/                    # Command-line applications
│   └── frr-tui/            # Your main TUI application
│       └── main.go         # Entry point
├── internal/               # Private application code
│   ├── ui/                 # UI components
│   │   ├── views/          # Different screens/views
│   │   ├── widgets/        # Reusable UI components
│   │   └── styles.go       # Common styling
│   ├── config/             # Configuration handling
│   ├── distro/             # Linux distro-specific code
│   │   ├── detector.go     # Distro detection logic
│   │   └── handlers/       # Distro-specific implementations
│   └── utils/              # Utility functions
├── pkg/                    # Public libraries that could be used by other projects
│   └── tuilib/             # Your reusable TUI components
├── assets/                 # Static assets (icons, etc.)
├── configs/                # Configuration files
├── docs/                   # Documentation
├── scripts/                # Build/deployment scripts
├── go.mod                  # Go module definition
└── README.md               # Project documentation
```

1. **Separation by functionality**: UI, distro handling, and configuration are clearly separated
2. **Easy extension**: Adding support for a new distro only requires adding a new handler in `internal/distro/handlers/`
3. **Private vs public code**: `internal/` keeps implementation details hidden while `pkg/` exposes reusable components
4. **Single responsibility**: Each package has a clear purpose
5. **Testability**: Components are modular and can be tested independently

This is the initial design of the code environment. 

## development Tools

### Linting and Code Quality

- **golangci-lint**: Linter for go
- **errcheck**: Ensures errors are correctly handled


### Testing and Coverage

- **go test**: Built-in testing framework

### Dependency Managemetn

- **go mod**: Standard Go module system

### Development Workflow

- **pre-commit hooks**: Run linters and tests before commits
- **Makefile**: Standardize common development commands
- **air**: Hot reloading during development

### IDE Integration

- **VS Code/codium** and **GoLand**: Both excellent Go development IDEs.

### CI/CD

- **GitHub Actions**: For this project we will use GitHub Actions for automated nightly builds and manual deployments.