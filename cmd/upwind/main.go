// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tektite-io/upwind-go-sdk/sdk"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Handle version command
	if command == "version" {
		fmt.Printf("upwind-cli version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Handle help command
	if command == "help" || command == "--help" || command == "-h" {
		printUsage()
		os.Exit(0)
	}

	// Create client
	client, err := createClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Enable logging if verbose flag is set
	if hasFlag("--verbose") || hasFlag("-v") {
		client.EnableLogging()
	}

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	// Execute command
	if err := executeCommand(ctx, client, command, os.Args[2:]); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func createClient() (*sdk.Client, error) {
	// Check if config file is specified
	configFile := getFlagValue("--config")
	if configFile != "" {
		return sdk.NewClientFromFile(configFile)
	}

	// Otherwise, use environment variables
	return sdk.NewClientFromEnv()
}

func executeCommand(ctx context.Context, client *sdk.Client, command string, args []string) error {
	switch command {
	case "vulnerability-findings", "vulns":
		return handleVulnerabilityFindings(ctx, client, args)
	case "configuration-findings", "config-findings":
		return handleConfigurationFindings(ctx, client, args)
	case "threat-detections", "threats":
		return handleThreatDetections(ctx, client, args)
	case "threat-events":
		return handleThreatEvents(ctx, client, args)
	case "threat-policies":
		return handleThreatPolicies(ctx, client, args)
	case "api-endpoints":
		return handleApiEndpoints(ctx, client, args)
	case "sbom-packages", "packages":
		return handleSbomPackages(ctx, client, args)
	case "workflows":
		return handleWorkflows(ctx, client, args)
	case "integration-webhooks", "webhooks":
		return handleIntegrationWebhooks(ctx, client, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleVulnerabilityFindings(ctx context.Context, client *sdk.Client, args []string) error {
	if len(args) > 0 && args[0] == "get" {
		if len(args) < 2 {
			return fmt.Errorf("usage: vulnerability-findings get <finding-id>")
		}
		finding, err := client.GetVulnerabilityFinding(ctx, args[1])
		if err != nil {
			return err
		}
		return printJSON(finding)
	}

	// List findings
	query := &sdk.VulnerabilityFindingsQuery{}
	if severity := getFlagValue("--severity"); severity != "" {
		query.Severity = severity
	}
	if imageName := getFlagValue("--image-name"); imageName != "" {
		query.ImageName = imageName
	}
	if hasFlag("--in-use") {
		inUse := true
		query.InUse = &inUse
	}
	if hasFlag("--exploitable") {
		exploitable := true
		query.Exploitable = &exploitable
	}

	findingsCh, errCh := client.ListVulnerabilityFindings(ctx, query)
	return streamAndPrintJSON(ctx, findingsCh, errCh)
}

func handleConfigurationFindings(ctx context.Context, client *sdk.Client, args []string) error {
	if len(args) > 0 && args[0] == "get" {
		if len(args) < 2 {
			return fmt.Errorf("usage: configuration-findings get <finding-id>")
		}
		finding, err := client.GetConfigurationFinding(ctx, args[1], hasFlag("--include-tags"))
		if err != nil {
			return err
		}
		return printJSON(finding)
	}

	// List findings
	query := &sdk.ConfigurationFindingsQuery{}
	if severity := getFlagValue("--severity"); severity != "" {
		query.Severity = severity
	}
	if status := getFlagValue("--status"); status != "" {
		query.Status = status
	}
	if frameworkID := getFlagValue("--framework-id"); frameworkID != "" {
		query.FrameworkID = frameworkID
	}

	findingsCh, errCh := client.ListConfigurationFindings(ctx, query)
	return streamAndPrintJSON(ctx, findingsCh, errCh)
}

func handleThreatDetections(ctx context.Context, client *sdk.Client, args []string) error {
	if len(args) > 0 && args[0] == "get" {
		if len(args) < 2 {
			return fmt.Errorf("usage: threat-detections get <detection-id>")
		}
		detection, err := client.GetThreatDetection(ctx, args[1])
		if err != nil {
			return err
		}
		return printJSON(detection)
	}

	if len(args) > 0 && args[0] == "archive" {
		if len(args) < 2 {
			return fmt.Errorf("usage: threat-detections archive <detection-id>")
		}
		detection, err := client.ArchiveThreatDetection(ctx, args[1])
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Detection archived successfully")
		return printJSON(detection)
	}

	// List detections
	query := &sdk.ThreatDetectionsQuery{}
	if severity := getFlagValue("--severity"); severity != "" {
		query.Severity = severity
	}
	if category := getFlagValue("--category"); category != "" {
		query.Category = category
	}
	if detectionType := getFlagValue("--type"); detectionType != "" {
		query.Type = detectionType
	}

	detections, err := client.ListThreatDetections(ctx, query)
	if err != nil {
		return err
	}
	return printJSON(detections)
}

func handleThreatEvents(ctx context.Context, client *sdk.Client, args []string) error {
	query := &sdk.ThreatEventsQuery{}
	if severity := getFlagValue("--severity"); severity != "" {
		query.Severity = severity
	}
	if category := getFlagValue("--category"); category != "" {
		query.Category = category
	}

	events, err := client.ListThreatEvents(ctx, query)
	if err != nil {
		return err
	}
	return printJSON(events)
}

func handleThreatPolicies(ctx context.Context, client *sdk.Client, args []string) error {
	managedBy := getFlagValue("--managed-by")
	policies, err := client.ListThreatPolicies(ctx, managedBy)
	if err != nil {
		return err
	}
	return printJSON(policies)
}

func handleApiEndpoints(ctx context.Context, client *sdk.Client, args []string) error {
	query := &sdk.ApiEndpointsQuery{}
	if method := getFlagValue("--method"); method != "" {
		query.Method = method
	}
	if domain := getFlagValue("--domain"); domain != "" {
		query.Domain = domain
	}
	if authState := getFlagValue("--auth-state"); authState != "" {
		query.AuthenticationState = authState
	}

	endpointsCh, errCh := client.ListApiEndpoints(ctx, query)
	return streamAndPrintJSON(ctx, endpointsCh, errCh)
}

func handleSbomPackages(ctx context.Context, client *sdk.Client, args []string) error {
	if len(args) > 0 && args[0] == "get" {
		if len(args) < 3 {
			return fmt.Errorf("usage: sbom-packages get <package-name> <version>")
		}
		pkg, err := client.GetSbomPackageDetails(ctx, args[1], args[2])
		if err != nil {
			return err
		}
		return printJSON(pkg)
	}

	// List packages
	query := &sdk.SbomPackagesQuery{}
	if packageName := getFlagValue("--package-name"); packageName != "" {
		query.PackageName = packageName
	}
	if framework := getFlagValue("--framework"); framework != "" {
		query.Framework = framework
	}

	packages, err := client.ListSbomPackages(ctx, query)
	if err != nil {
		return err
	}
	return printJSON(packages)
}

func handleWorkflows(ctx context.Context, client *sdk.Client, args []string) error {
	if len(args) > 0 && args[0] == "get" {
		if len(args) < 2 {
			return fmt.Errorf("usage: workflows get <workflow-id>")
		}
		workflow, err := client.GetWorkflow(ctx, args[1])
		if err != nil {
			return err
		}
		return printJSON(workflow)
	}

	// List workflows
	workflows, err := client.ListWorkflows(ctx)
	if err != nil {
		return err
	}
	return printJSON(workflows)
}

func handleIntegrationWebhooks(ctx context.Context, client *sdk.Client, args []string) error {
	vendor := getFlagValue("--vendor")
	webhooks, err := client.ListIntegrationWebhooks(ctx, vendor)
	if err != nil {
		return err
	}
	return printJSON(webhooks)
}

// Helper functions

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func streamAndPrintJSON[T any](ctx context.Context, itemsCh <-chan T, errCh <-chan error) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	// Start JSON array
	fmt.Println("[")
	first := true

	for {
		select {
		case item, ok := <-itemsCh:
			if !ok {
				// Channel closed
				fmt.Println("]")
				// Check for errors
				select {
				case err := <-errCh:
					return err
				default:
					return nil
				}
			}
			if !first {
				fmt.Println(",")
			}
			first = false
			if err := encoder.Encode(item); err != nil {
				return fmt.Errorf("encoding JSON: %w", err)
			}
		case err := <-errCh:
			if err != nil {
				fmt.Println("]")
				return err
			}
		case <-ctx.Done():
			fmt.Println("]")
			return ctx.Err()
		}
	}
}

func hasFlag(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

func getFlagValue(flag string) string {
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}
	return ""
}

func printUsage() {
	fmt.Printf(`Upwind CLI %s - Command-line interface for the Upwind Security API

USAGE:
    upwind <command> [subcommand] [options]

COMMANDS:
    vulnerability-findings, vulns          List or get vulnerability findings
    configuration-findings, config-findings List or get configuration findings
    threat-detections, threats             List, get, or archive threat detections
    threat-events                          List threat events
    threat-policies                        List threat policies
    api-endpoints                          List API security endpoints
    sbom-packages, packages                List or get SBOM packages
    workflows                              List or get workflows
    integration-webhooks, webhooks         List integration webhooks
    version                                Show version information
    help                                   Show this help message

SUBCOMMANDS:
    get <id>                              Get a specific resource by ID
    archive <id>                          Archive a resource (threat detections only)

GLOBAL OPTIONS:
    --config <file>                       Path to config file (JSON format)
    --verbose, -v                         Enable verbose logging
    -h, --help                            Show help message

FILTER OPTIONS (varies by command):
    --severity <severity>                 Filter by severity (LOW, MEDIUM, HIGH, CRITICAL)
    --status <status>                     Filter by status
    --category <category>                 Filter by category
    --type <type>                         Filter by type
    --image-name <name>                   Filter by image name
    --framework-id <id>                   Filter by framework ID
    --managed-by <entity>                 Filter by managed-by entity
    --method <method>                     Filter by HTTP method
    --domain <domain>                     Filter by domain
    --auth-state <state>                  Filter by authentication state
    --package-name <name>                 Filter by package name
    --framework <framework>               Filter by framework
    --vendor <vendor>                     Filter by vendor
    --in-use                              Filter for packages in use
    --exploitable                         Filter for exploitable vulnerabilities
    --include-tags                        Include cloud account tags in response

ENVIRONMENT VARIABLES:
    UPWIND_CLIENT_ID                      OAuth2 client ID (required)
    UPWIND_CLIENT_SECRET                  OAuth2 client secret (required)
    UPWIND_ORGANIZATION_ID                Organization ID (required)
    UPWIND_REGION                         API region: US, EU, or ME (default: US)
    UPWIND_BASE_URL                       Custom base URL (optional)
    UPWIND_TOKEN_URL                      Custom token URL (optional)
    UPWIND_MAX_RETRIES                    Maximum retry attempts (default: 3)
    UPWIND_MAX_CONCURRENCY                Maximum concurrent requests (default: 10)
    UPWIND_PAGE_SIZE                      Default page size (default: 100)
    UPWIND_RATE_LIMIT                     Requests per second limit (default: 10)

EXAMPLES:
    # List all vulnerability findings with high severity
    upwind vulnerability-findings --severity HIGH

    # Get a specific configuration finding
    upwind configuration-findings get finding_123456

    # Archive a threat detection
    upwind threat-detections archive detection_789

    # List API endpoints for a specific domain
    upwind api-endpoints --domain api.example.com

    # Get SBOM package details
    upwind sbom-packages get lodash 4.17.21

    # Use a config file instead of environment variables
    upwind --config config.json vulnerability-findings

For more information, visit: https://docs.upwind.io
`, sdk.Version)
}
