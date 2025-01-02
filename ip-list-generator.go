package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config holds all program configuration parameters
type Config struct {
	cidr      string // CIDR range for IP generation
	outputDir string // Directory to save output file
	filename  string // Custom filename (optional)
}

// main is the entry point of the application
func main() {
	// Parse command line flags and get configuration
	config := parseFlags()

	// Generate IPs and handle any errors
	if err := generateIPs(config); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags processes command line arguments and returns a Config struct
func parseFlags() *Config {
	config := &Config{}

	// Define command line flags
	flag.StringVar(&config.cidr, "cidr", "", "CIDR range (e.g., 192.168.1.0/24)")
	flag.StringVar(&config.outputDir, "output", "", "Output directory path")
	flag.StringVar(&config.filename, "filename", "", "Custom filename (optional)")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if config.cidr == "" {
		fmt.Println("Error: CIDR range is required")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return config
}

// generateIPs handles the IP generation and file writing process
func generateIPs(config *Config) error {
	// Validate and parse CIDR notation
	ip, ipnet, err := net.ParseCIDR(config.cidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR format: %v", err)
	}

	// Set default output directory if not specified
	if config.outputDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		config.outputDir = currentDir
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Generate default filename if not provided
	if config.filename == "" {
		timestamp := time.Now().Format("20060102_150405")
		sanitizedCIDR := strings.Replace(config.cidr, "/", "_", -1)
		sanitizedCIDR = strings.Replace(sanitizedCIDR, ".", "-", -1)
		config.filename = fmt.Sprintf("ip_list_%s_%s.txt", sanitizedCIDR, timestamp)
	}

	// Ensure filename has .txt extension
	if !strings.HasSuffix(config.filename, ".txt") {
		config.filename += ".txt"
	}

	// Construct full file path
	filepath := filepath.Join(config.outputDir, config.filename)

	// Create and open output file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create buffered writer for better performance
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Initialize progress tracking
	count := 0
	startTime := time.Now()

	// Generate and write IPs
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if _, err := writer.WriteString(ip.String() + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
		count++

		// Show progress for large ranges
		if count%10000 == 0 {
			fmt.Printf("Generated %d IPs...\n", count)
		}
	}

	// Calculate execution time
	duration := time.Since(startTime)

	// Print summary
	fmt.Printf("\nExecution Summary:\n")
	fmt.Printf("----------------\n")
	fmt.Printf("CIDR Range: %s\n", config.cidr)
	fmt.Printf("Total IPs Generated: %d\n", count)
	fmt.Printf("Time Taken: %v\n", duration)
	fmt.Printf("Output File: %s\n", filepath)
	fmt.Printf("Average Speed: %.2f IPs/second\n", float64(count)/duration.Seconds())

	return nil
}

// inc increments an IP address by one
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// validatePath checks if a path is valid and accessible
func validatePath(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error accessing path: %v", err)
	}

	// Check if path is writable
	tmpFile := filepath.Join(path, ".tmp")
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("path is not writable: %s", path)
	}
	f.Close()
	os.Remove(tmpFile)

	return nil
}
