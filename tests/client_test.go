package checkhost_test

import (
	"log"
	"testing"
	"time"

	checkhost "github.com/check-hostcc/check-host-api-go"
)

func setUp() *checkhost.CheckHost {
	// Initialize without API key for public tests
	return checkhost.NewClient("")
}

func throttle() {
	log.Println("Sleeping 5 seconds to respect rate limits...")
	time.Sleep(5 * time.Second)
}

func TestUtilities(t *testing.T) {
	client := setUp()

	t.Run("MyIP", func(t *testing.T) {
		ip, err := client.MyIP()
		if err != nil {
			t.Fatalf("MyIP failed: %v", err)
		}
		if ip == "" {
			t.Error("Expected an IP, got empty string")
		}
		log.Printf("MyIP: %s\n", ip)
	})

	throttle()

	t.Run("Locations", func(t *testing.T) {
		nodes, err := client.Locations()
		if err != nil {
			t.Fatalf("Locations failed: %v", err)
		}
		if len(nodes) == 0 {
			t.Error("Expected nodes, got none")
		}
		log.Printf("Successfully fetched %d nodes\n", len(nodes))
	})

	throttle()

	t.Run("Info", func(t *testing.T) {
		info, err := client.Info("check-host.cc")
		if err != nil {
			t.Fatalf("Info failed: %v", err)
		}
		if info.IP == "" {
			t.Error("Expected IP in info response")
		}
		log.Printf("Info country: %s\n", info.Country)
	})

	throttle()

	t.Run("Whois", func(t *testing.T) {
		whoisRaw, err := client.Whois("check-host.cc")
		if err != nil {
			t.Fatalf("Whois failed: %v", err)
		}
		if whoisRaw == "" {
			t.Error("Expected whois output")
		}
		log.Printf("Whois length: %d chars\n", len(whoisRaw))
	})
}

func TestActiveMonitoring(t *testing.T) {
	client := setUp()
	target := "8.8.8.8"
	domainTarget := "check-host.cc"

	throttle()

	t.Run("Ping", func(t *testing.T) {
		res, err := client.Ping(target, &checkhost.MonitoringRequest{
			Region: []string{"US"},
		})
		if err != nil {
			t.Fatalf("Ping failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("Ping UUID: %s\n", res.UUID)

		// Test Report on the generated UUID immediately just to verify the endpoint is functioning
		throttle()
		report, err := client.Report(res.UUID)
		if err != nil {
			t.Fatalf("Report failed for UUID %s: %v", res.UUID, err)
		}
		log.Printf("Report fetched %d keys for %s", len(report), res.UUID)
	})

	throttle()

	t.Run("DNS", func(t *testing.T) {
		res, err := client.DNS(domainTarget, &checkhost.DNSTargetRequest{
			QueryMethod: "TXT",
			Region:      []string{"DE"},
		})
		if err != nil {
			t.Fatalf("DNS failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("DNS UUID: %s\n", res.UUID)
	})

	throttle()

	t.Run("TCP", func(t *testing.T) {
		res, err := client.TCP(target, 53, &checkhost.TCPMonitoringRequest{})
		if err != nil {
			t.Fatalf("TCP failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("TCP UUID: %s\n", res.UUID)
	})

	throttle()

	t.Run("UDP", func(t *testing.T) {
		res, err := client.UDP(target, 53, &checkhost.UDPMonitoringRequest{})
		if err != nil {
			t.Fatalf("UDP failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("UDP UUID: %s\n", res.UUID)
	})

	throttle()

	t.Run("HTTP", func(t *testing.T) {
		res, err := client.Http("https://check-host.cc", &checkhost.MonitoringRequest{})
		if err != nil {
			t.Fatalf("HTTP failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("HTTP UUID: %s\n", res.UUID)
	})

	throttle()

	t.Run("MTR", func(t *testing.T) {
		res, err := client.MTR(target, &checkhost.MTRMonitoringRequest{
			Region: []string{"US"},
		})
		if err != nil {
			t.Fatalf("MTR failed: %v", err)
		}
		if res.UUID == "" {
			t.Error("Expected UUID")
		}
		log.Printf("MTR UUID: %s\n", res.UUID)
	})
}
