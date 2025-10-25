package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"portscanner/internal/scanner"
)

func main() {
	var (
		host    = flag.String("host", "", "Target host or IP address (required)")
		ports   = flag.String("p", "", "Port range (e.g., '80', '1-1000', '22,80,443') (required)")
		timeout = flag.Int("t", 500, "Timeout in milliseconds")
		workers = flag.Int("workers", 0, "Limit of concurrent goroutines (0 = unlimited)")
		banner  = flag.Bool("banner", false, "Collect banners from open ports")
	)
	flag.Parse()

	if *host == "" || *ports == "" {
		fmt.Println("Error: -host and -p parameters are required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse port range
	portRange, err := scanner.ParsePortRange(*ports)
	if err != nil {
		log.Fatalf("Error parsing port range: %v", err)
	}

	// Create scanner config
	config := &scanner.Config{
		Host:      *host,
		PortRange: portRange,
		Timeout:   *timeout,
		Workers:   *workers,
		GetBanner: *banner,
	}

	// Run scan
	results := scanner.Scan(config)

	// Print results
	for _, result := range results {
		if result.Open {
			if *banner && result.Banner != "" {
				fmt.Printf("%d/tcp open %s %s\n", result.Port, result.Service, result.Banner)
			} else {
				fmt.Printf("%d/tcp open\n", result.Port)
			}
		}
	}
}
