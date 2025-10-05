package main

import (
	"fmt"
	"os"
	"strings"

	"jsonrpc/internal/executor"
	"jsonrpc/internal/output"
	"jsonrpc/internal/parser"
	"jsonrpc/pkg/types"

	"github.com/spf13/cobra"
)

var (
	// Version information (set by GoReleaser)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	jsonOutput bool
	detailed   bool

	// Run command flags
	urlFlag     string
	headerFlags []string
	configFlag  string
	timeoutFlag int
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc-cli",
		Short: "Execute JSON-RPC requests defined in HCL configuration files",
		Long: `rpc-cli is a CLI tool that reads JSON-RPC request definitions from HCL files
and executes them with flexible configuration management.`,
		Version: version,
	}

	cmd.AddCommand(
		lsCmd(),
		runCmd(),
		validateCmd(),
		versionCmd(),
	)

	return cmd
}

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls <file> [request_names...]",
		Short: "List requests from HCL file",
		Long: `List all requests or specific requests from an HCL file.
With no request names, lists all requests.
With request names, lists only specified requests.`,
		Args: cobra.MinimumNArgs(1),
		RunE: runListCommand,
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed information")

	return cmd
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <file> [request_names...]",
		Short: "Execute requests",
		Long: `Execute all requests or specific requests from an HCL file.
With no request names, executes all requests.
With request names, executes only specified requests.`,
		Args: cobra.MinimumNArgs(1),
		RunE: runExecuteCommand,
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.Flags().StringVar(&urlFlag, "url", "", "Override URL for requests")
	cmd.Flags().StringArrayVar(&headerFlags, "header", []string{}, "Override headers (can be repeated)")
	cmd.Flags().StringVar(&configFlag, "config", "", "Use specific config profile")
	cmd.Flags().IntVar(&timeoutFlag, "timeout", 0, "Override timeout in seconds")

	return cmd
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate HCL syntax",
		Long:  `Validate HCL file syntax, check required fields, and verify config references.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runValidateCommand,
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rpc-cli version %s\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built:  %s\n", date)
		},
	}
}

func runListCommand(cmd *cobra.Command, args []string) error {
	filename := args[0]
	requestNames := args[1:]

	// Parse HCL file
	p := parser.New()
	hclFile, err := p.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	// Filter requests if names specified
	requestsToShow, err := filterRequests(hclFile, requestNames)
	if err != nil {
		return err
	}

	// Format and output results
	formatter := output.New()

	if jsonOutput {
		return formatter.FormatRequestJSON(requestsToShow)
	} else if detailed {
		formatter.FormatRequestDetailed(hclFile, requestsToShow, nil)
	} else {
		formatter.FormatRequestList(hclFile, requestsToShow, nil)
	}

	return nil
}

func runExecuteCommand(cmd *cobra.Command, args []string) error {
	filename := args[0]
	requestNames := args[1:]

	// Parse HCL file
	p := parser.New()
	hclFile, err := p.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	// Validate HCL file
	validator := parser.NewValidator()
	if err := validator.Validate(hclFile); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Filter requests if names specified
	requestsToRun, err := filterRequests(hclFile, requestNames)
	if err != nil {
		return err
	}

	// Build CLI overrides
	overrides, err := buildCLIOverrides()
	if err != nil {
		return err
	}

	// Execute requests
	exec := executor.New()
	results, err := exec.ExecuteAll(hclFile, requestsToRun, overrides)
	if err != nil {
		return fmt.Errorf("failed to execute requests: %w", err)
	}

	// Format and output results
	formatter := output.New()
	formatter.FormatExecutionResults(results, jsonOutput)

	// Exit with error code if any request failed
	for _, result := range results {
		if !result.IsSuccess() {
			os.Exit(1)
		}
	}

	return nil
}

func runValidateCommand(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Parse HCL file
	p := parser.New()
	hclFile, err := p.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	// Validate HCL file
	validator := parser.NewValidator()
	if err := validator.Validate(hclFile); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("✓ File '%s' is valid\n", filename)
	fmt.Printf("  - %d config(s) found\n", len(hclFile.Configs))
	fmt.Printf("  - %d request(s) found\n", len(hclFile.Requests))

	return nil
}

// filterRequests filters requests by name if specified
func filterRequests(hclFile *types.HCLFile, requestNames []string) ([]*types.Request, error) {
	if len(requestNames) == 0 {
		return hclFile.Requests, nil
	}

	requestMap := make(map[string]*types.Request)
	for _, req := range hclFile.Requests {
		requestMap[req.Name] = req
	}

	var filtered []*types.Request
	for _, name := range requestNames {
		req, exists := requestMap[name]
		if !exists {
			return nil, fmt.Errorf("request '%s' not found in file", name)
		}
		filtered = append(filtered, req)
	}

	return filtered, nil
}

// buildCLIOverrides builds CLI overrides from flags
func buildCLIOverrides() (*types.CLIOverrides, error) {
	overrides := types.NewCLIOverrides()
	overrides.URL = urlFlag
	overrides.Config = configFlag
	overrides.Timeout = timeoutFlag

	// Parse header flags
	for _, header := range headerFlags {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s (expected 'Key: Value')", header)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		overrides.Headers[key] = value
	}

	return overrides, nil
}
