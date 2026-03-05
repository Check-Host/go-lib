package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	checkhost "github.com/check-hostcc/check-host-api-go"
)

func main() {
	// Initialize the client. The API Key is optional.
	// Without an API key, standard public rate limits apply.
	// client := checkhost.NewClient("YOUR_API_KEY_HERE")
	client := checkhost.NewClient("")

	fmt.Println("=== Check-Host API Go Example ===")

	// Fetch all locations
	fmt.Println("\n1. Fetching available nodes...")
	nodes, err := client.Locations()
	if err != nil {
		log.Fatalf("Error fetching locations: %v", err)
	}
	fmt.Printf("Successfully retrieved %d global nodes.\n", len(nodes))

	// Get My IP
	fmt.Println("\n2. Getting my public IP...")
	ip, err := client.MyIP()
	if err != nil {
		log.Fatalf("Error getting IP: %v", err)
	}
	fmt.Printf("Your IP is: %s\n", ip)

	// Sleep to respect public rate limiting as we don't have an API key in this example
	fmt.Println("\nSleeping 5 seconds to respect rate limits...")
	time.Sleep(5 * time.Second)

	// Perform a Ping Check
	target := "8.8.8.8"
	fmt.Printf("\n3. Starting Ping Check for %s (EU Region)...\n", target)

	pingRes, err := client.Ping(target, &checkhost.MonitoringRequest{
		Region: []string{"EU"},
	})
	if err != nil {
		log.Fatalf("Error starting check: %v", err)
	}

	fmt.Printf("Check started successfully! Task UUID: %s\n", pingRes.UUID)
	fmt.Printf("You can view real-time results via browser: %s\n", pingRes.ReportURL)

	// Wait 2 seconds before polling results
	fmt.Println("Polling report in 2 seconds...")
	time.Sleep(2 * time.Second)

	report, err := client.Report(pingRes.UUID)
	if err != nil {
		log.Fatalf("Error fetching report: %v", err)
	}

	// Unmarshal and pretty-print report json
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Printf("Preliminary Report Results: \n%s\n", string(reportJSON))
}
