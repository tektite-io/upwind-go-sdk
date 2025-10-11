# Upwind Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/tektite-io/upwind-go-sdk.svg)](https://pkg.go.dev/github.com/tektite-io/upwind-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/tektite-io/upwind-go-sdk)](https://goreportcard.com/report/github.com/tektite-io/upwind-go-sdk)

A production-ready Go SDK and CLI for the [Upwind Security](https://upwind.io) API. Designed for efficiency, suitable for Docker containers, Steampipe plugins, and CloudQuery plugins.

## Features

- 🚀 **Production Ready** - Memory-efficient (<4GB RAM), low CPU usage (~0.5 CPU)
- 🔄 **Streaming Support** - Process large datasets without loading everything into memory
- 🔐 **OAuth2 Authentication** - Automatic token management and refresh
- ⚡ **Rate Limiting** - Built-in rate limiting with configurable requests per second
- 🔁 **Smart Retries** - Exponential backoff with configurable retry attempts
- 🎯 **Context Aware** - All operations respect context cancellation
- 🛠️ **CLI & SDK** - Use as a library or standalone command-line tool
- 🐳 **Docker Support** - Minimal container image (~10MB, built from scratch)

## Installation

### As a Go Module

```bash
go get github.com/tektite-io/upwind-go-sdk
```

### CLI Binary

**Using Go:**
```bash
go install github.com/tektite-io/upwind-go-sdk/cmd/upwind@v1.0.0
```

**From Source:**
```bash
git clone https://github.com/tektite-io/upwind-go-sdk.git
cd upwind-go-sdk
make build
./build/upwind version
```

**Using Docker:**
```bash
docker pull ghcr.io/tektite-io/upwind-go-sdk:1.0.0
```

## Quick Start

### Using the SDK

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/tektite-io/upwind-go-sdk/sdk"
)

func main() {
    // Create client from environment variables
    client, err := sdk.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Stream vulnerability findings
    opts := sdk.ListVulnerabilityFindingsOptions{
        Severity: "CRITICAL",
    }

    resultsCh, errCh := client.ListVulnerabilityFindings(ctx, opts)

    for result := range resultsCh {
        fmt.Printf("Finding: %s - %s\n", result.ID, result.Title)
    }

    if err := <-errCh; err != nil {
        log.Fatal(err)
    }
}
```

### Using the CLI

```bash
# Set up credentials
export UPWIND_CLIENT_ID="your-client-id"
export UPWIND_CLIENT_SECRET="your-client-secret"
export UPWIND_ORGANIZATION_ID="your-org-id"
export UPWIND_REGION="US"  # US, EU, or ME

# List critical vulnerability findings
upwind vulnerability-findings --severity CRITICAL

# Get specific finding
upwind vulnerability-findings get <finding-id>

# List API endpoints
upwind api-endpoints --domain example.com

# View threat detections
upwind threat-detections --status OPEN
```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `UPWIND_CLIENT_ID` | OAuth2 client ID | Yes | - |
| `UPWIND_CLIENT_SECRET` | OAuth2 client secret | Yes | - |
| `UPWIND_ORGANIZATION_ID` | Organization ID | Yes | - |
| `UPWIND_REGION` | API region (US, EU, ME) | No | US |
| `UPWIND_MAX_RETRIES` | Maximum retry attempts | No | 3 |
| `UPWIND_MAX_CONCURRENCY` | Max concurrent requests | No | 10 |
| `UPWIND_PAGE_SIZE` | Items per page | No | 100 |
| `UPWIND_RATE_LIMIT_PER_SECOND` | Rate limit (req/sec) | No | 10 |

### Configuration File

Create a `config.json` file:

```json
{
  "client_id": "your_client_id",
  "client_secret": "your_client_secret",
  "organization_id": "your_org_id",
  "region": "US",
  "max_retries": 3,
  "max_concurrency": 10,
  "page_size": 100,
  "rate_limit_per_second": 10
}
```

**Load from file:**
```go
client, err := sdk.NewClientFromFile("config.json")
```

Or with the CLI:
```bash
upwind --config config.json vulnerability-findings
```

## SDK Usage Examples

### Vulnerability Findings

```go
// List with filters
opts := sdk.ListVulnerabilityFindingsOptions{
    Severity:     "HIGH",
    Status:       "OPEN",
    Exploitable:  sdk.Bool(true),
    ImageName:    "nginx:latest",
}

resultsCh, errCh := client.ListVulnerabilityFindings(ctx, opts)

// Collect all results
findings, err := sdk.CollectAll(resultsCh, errCh)
if err != nil {
    log.Fatal(err)
}

// Get specific finding
finding, err := client.GetVulnerabilityFinding(ctx, "finding-id")
```

### Configuration Findings

```go
opts := sdk.ListConfigurationFindingsOptions{
    Severity:    "CRITICAL",
    Category:    "Security",
    FrameworkID: "cis-aws",
}

resultsCh, errCh := client.ListConfigurationFindings(ctx, opts)

for result := range resultsCh {
    fmt.Printf("Config Issue: %s\n", result.Title)
}
```

### Threat Detections

```go
// List threats
opts := sdk.ListThreatDetectionsOptions{
    Status:   "OPEN",
    Severity: "HIGH",
}

threatsCh, errCh := client.ListThreatDetections(ctx, opts)

// Archive a threat
err := client.ArchiveThreatDetection(ctx, "threat-id")
```

### API Security Endpoints

```go
opts := sdk.ListAPIEndpointsOptions{
    Domain:    "api.example.com",
    Method:    "GET",
    AuthState: "AUTHENTICATED",
}

endpointsCh, errCh := client.ListAPIEndpoints(ctx, opts)
```

### SBOM Packages

```go
opts := sdk.ListSBOMPackagesOptions{
    PackageName: "openssl",
    InUse:       sdk.Bool(true),
}

packagesCh, errCh := client.ListSBOMPackages(ctx, opts)
```

### Workflows

```go
// List workflows
workflowsCh, errCh := client.ListWorkflows(ctx, sdk.ListWorkflowsOptions{})

// Get workflow
workflow, err := client.GetWorkflow(ctx, "workflow-id")

// List integration webhooks
webhooksCh, errCh := client.ListIntegrationWebhooks(ctx, sdk.ListIntegrationWebhooksOptions{})
```

## CLI Commands

### Vulnerability Findings

```bash
# List all
upwind vulnerability-findings

# Filter by severity
upwind vulns --severity CRITICAL

# Filter exploitable
upwind vulns --exploitable

# Get specific finding
upwind vulns get <id>
```

### Configuration Findings

```bash
# List all
upwind configuration-findings

# Filter by category and severity
upwind config-findings --category Security --severity HIGH

# Get specific finding
upwind config-findings get <id>
```

### Threat Detections

```bash
# List threats
upwind threat-detections --status OPEN

# Get specific threat
upwind threats get <id>

# Archive threat
upwind threats archive <id>

# List threat events
upwind threat-events

# List threat policies
upwind threat-policies
```

### API Security

```bash
# List API endpoints
upwind api-endpoints

# Filter by domain
upwind api-endpoints --domain api.example.com

# Filter by HTTP method
upwind api-endpoints --method POST

# Filter authenticated endpoints
upwind api-endpoints --auth-state AUTHENTICATED
```

### SBOM Packages

```bash
# List all packages
upwind sbom-packages

# Filter by name
upwind packages --package-name openssl

# Filter in-use packages
upwind packages --in-use
```

### Workflows

```bash
# List workflows
upwind workflows

# Get specific workflow
upwind workflows get <id>

# List webhooks
upwind integration-webhooks
```

## Docker Usage

### Build Image

```bash
make docker-build
```

### Run CLI in Docker

```bash
docker run --rm \
  -e UPWIND_CLIENT_ID=$UPWIND_CLIENT_ID \
  -e UPWIND_CLIENT_SECRET=$UPWIND_CLIENT_SECRET \
  -e UPWIND_ORGANIZATION_ID=$UPWIND_ORGANIZATION_ID \
  -e UPWIND_REGION=US \
  upwind-go-sdk:1.0.0 vulnerability-findings --severity CRITICAL
```

Or use the Makefile:

```bash
make docker-run
```

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile)

### Building

```bash
# Install dependencies
make deps

# Build binary
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench
```

## API Reference

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/tektite-io/upwind-go-sdk).

### Core Types

- `Client` - Main SDK client
- `Config` - Configuration options
- `VulnerabilityFinding` - Vulnerability finding details
- `ConfigurationFinding` - Configuration issue details
- `ThreatDetection` - Threat detection details
- `APIEndpoint` - API security endpoint details
- `SBOMPackage` - SBOM package information
- `Workflow` - Workflow definition

### Helper Functions

- `CollectAll(resultsCh, errCh)` - Collect all streaming results into a slice
- `sdk.Bool(v)` - Helper for creating boolean pointers
- `sdk.Int(v)` - Helper for creating integer pointers
- `sdk.String(v)` - Helper for creating string pointers

## Performance Characteristics

- **Memory**: <4GB RAM for typical workloads
- **CPU**: ~0.5 CPU average usage
- **Container Size**: ~10MB Docker image (scratch-based)
- **Concurrency**: Configurable with semaphore-based limiting
- **Rate Limiting**: Built-in with configurable req/sec

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Mozilla Public License 2.0 - see the LICENSE file for details.

## Support

- 📧 Email: info@tektite.io
- 🐛 Issues: [GitHub Issues](https://github.com/tektite-io/upwind-go-sdk/issues)
- 📚 Documentation: [pkg.go.dev](https://pkg.go.dev/github.com/tektite-io/upwind-go-sdk)

## Acknowledgments

Built with ❤️ by [Tektite IO](https://tektite.io) for the [Upwind Security](https://upwind.io) platform.

