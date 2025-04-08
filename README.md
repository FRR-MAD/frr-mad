# frr-tui

## Usage

## Installation

- foobar

    ```


## Project Structure

```
root/
├── frontend/               # 
├── backend/                # 
│   ├── analytics/          # 
│   └── collector/          # 
└── README.md               # Project documentation
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