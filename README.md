# Check-Host API Go Library

A lightweight, lightning-fast, and feature-complete Go wrapper for the [Check-Host.cc](https://check-host.cc) API. Full API reference: [check-host.cc/docs](https://check-host.cc/docs). A bundled OpenAPI 3.0.3 / Swagger spec ships at [`swagger.yaml`](./swagger.yaml) for codegen / offline browsing.

Seamlessly integrate global network diagnostics into your backend. Perform remote Ping, MTR, DNS, HTTP, TCP and UDP checks from multiple worldwide locations—straight from your Go application. Checks from 60+ locations worldwide.

## Features

- **Zero Dependencies:** Built purely on the native Go `net/http` standard library. Zero package bloat.
- **Bulletproof Payloads:** Strictly utilizes POST requests for all active monitoring endpoints. This completely eliminates nasty URL-encoding issues with complex hostnames or custom UDP payloads.
- **Modern & Clean:** Written idiomatically with clear configuration structures and typed responses.
- **Smart Authentication:** API Key auto-injection. Configure your key once during client initialization, and the core SDK seamlessly handles all authentication payloads under the hood.

## Requirements

- **Go**: 1.18+

## Installation

Install the package directly using `go get`:
```bash
go get github.com/Check-Host/go-lib
```

## Quickstart

```go
package main

import (
	"fmt"
	"log"
	
	checkhost "github.com/Check-Host/go-lib"
)

func main() {
	// Initialize the client. The API Key is optional.
	// Without an API key, standard public rate limits apply.
	// client := checkhost.NewClient("YOUR_API_KEY_HERE")
	// Or leave empty: checkhost.NewClient("")
	client := checkhost.NewClient("")

	// Example: Retrieve all current nodes
	locations, err := client.Locations()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	fmt.Printf("Successfully retrieved %d global nodes.\n", len(locations))
}
```

---

## Complete API Reference & Examples

This library supports both minimal invocations and detailed, options-rich requests for every endpoint. All failures (network issues, API errors, rate limits) return standard `error` types encapsulating the actual check-host API response message.

### Common Options Used in Examples
Many endpoints accept a specific `Request` struct containing optional configuration fields:
- `Region`: Array of Nodes or ISO Country Codes (e.g. `[]string{"DE", "NL"}`) or Continents (e.g. `[]string{"EU"}`).
- `RepeatChecks`: Number of repeated probes to perform per node for higher accuracy (Live Check).
- `Timeout`: Connection timeout threshold in seconds. Supported by methods where a timeout is applicable (e.g., HTTP, TCP).

*Note: In Go, passing `nil` as the configuration object will automatically invoke the minimum configuration defaults required by the Check-Host API.*

---

### Information & Utilities

#### Get My IP
Returns the requesting client's public IPv4 or IPv6 address.
```go
ip, err := client.MyIP()
```

#### Get Locations
Fetches a dynamic list of all currently active monitoring nodes across the globe.
```go
nodes, err := client.Locations()
```

#### Host Info (GeoIP/ASN)
Retrieves detailed geolocation data, ISP information, and ASN details.
```go
// Minimal Example
info, err := client.Info("check-host.cc")
```

#### WHOIS Lookup
Performs a WHOIS registry lookup.
```go
// Minimal Example
whois, err := client.Whois("check-host.cc")
```

---

### Active Monitoring (POST Tasks)

Monitoring endpoints initiate tasks asynchronously and return a `CheckCreated` object containing an `UUID`. Use the `Report()` method (documented below) to fetch the actual results.

#### Ping
Dispatches ICMP echo requests to the target from global nodes.
```go
// Minimal Example
pingMin, err := client.Ping("8.8.8.8", nil)

// Max Example (With options)
pingMax, err := client.Ping("8.8.8.8", &checkhost.MonitoringRequest{
	Region:       []string{"DE", "NL"},
	RepeatChecks: 5,
	Timeout:      5,
})
```

#### DNS
Queries global nameservers for specific DNS records.
```go
// Minimal Example
dnsMin, err := client.DNS("check-host.cc", nil)

// Max Example (With options - TXT Record)
dnsMax, err := client.DNS("check-host.cc", &checkhost.DNSTargetRequest{
	QueryMethod: "TXT", // A, AAAA, MX, TXT, SRV, etc.
	Region:      []string{"US", "DE"},
})
```

#### TCP
Attempts to establish a 3-way TCP handshake on a specific destination port.
```go
// Minimal Example (Target, Port)
tcpMin, err := client.TCP("1.1.1.1", 443, nil)

// Max Example (With options)
tcpMax, err := client.TCP("1.1.1.1", 80, &checkhost.TCPMonitoringRequest{
	MonitoringRequest: checkhost.MonitoringRequest{
		Region:       []string{"DE", "NL"},
		RepeatChecks: 3,
		Timeout:      10,
	},
})
```

#### UDP
Sends UDP packets to a specified target and port.
```go
// Minimal Example (Target, Port)
udpMin, err := client.UDP("1.1.1.1", 53, nil)

// Max Example (With custom hex payload and options)
udpMax, err := client.UDP("1.1.1.1", 123, &checkhost.UDPMonitoringRequest{
	Payload: "0b", // NTP Request Hex
	MonitoringRequest: checkhost.MonitoringRequest{
		Region:       []string{"EU"},
		RepeatChecks: 2,
		Timeout:      5,
	},
})
```

#### HTTP
Executes an HTTP/HTTPS request to the target to measure TTFB and latency.
```go
// Minimal Example
httpMin, err := client.Http("https://check-host.cc", nil)

// Max Example (With options)
httpMax, err := client.Http("https://check-host.cc", &checkhost.MonitoringRequest{
	Region:       []string{"US", "DE"},
	RepeatChecks: 3,
	Timeout:      10,
})
```

#### MTR
Initiates an MTR (My Traceroute) diagnostic.
```go
// Minimal Example
mtrMin, err := client.MTR("1.1.1.1", nil)

// Max Example (With protocols, IP forced, and options)
mtrMax, err := client.MTR("1.1.1.1", &checkhost.MTRMonitoringRequest{
	RepeatChecks:   15,
	ForceIPVersion: 4,     // 4 or 6
	ForceProtocol:  "TCP", // default is ICMP
	Region:         []string{"DE", "US"},
})
```

---

### Fetching Results

#### Report
Fetches the compiled report and real-time statuses from a previously initiated monitoring check (Ping, TCP, HTTP, etc.) using its unique `UUID`. Wait 1-2 seconds after starting a check before polling. Longer checks with multiple repeats take one check per second and can be requested multiple times.
```go
// The check UUID is returned by any monitoring method above
taskUuid := "c0b4b0e3-aed7-4ae2-9f53-7bac879697cb"

// Fetch the result payload
report, err := client.Report(taskUuid)
```

## License

ISC License
