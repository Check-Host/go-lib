package checkhost

import (
	"fmt"
)

// MyIP returns the requesting client's public IPv4 or IPv6 address.
func (c *CheckHost) MyIP() (string, error) {
	var ip string
	err := c.doRequest("GET", "myip", nil, &ip)
	return ip, err
}

// Locations fetches a dynamic list of all currently active monitoring nodes across the globe.
func (c *CheckHost) Locations() (map[string]interface{}, error) {
	var nodes map[string]interface{}
	err := c.doRequest("GET", "locations", nil, &nodes)
	return nodes, err
}

// Info retrieves detailed geolocation data, ISP information, and ASN details.
func (c *CheckHost) Info(target string) (*MinResponseINFO, error) {
	req := TargetRequest{
		APIKey: c.APIKey,
		Target: target,
	}
	var info MinResponseINFO
	err := c.doRequest("POST", "info", req, &info)
	return &info, err
}

// Whois performs a WHOIS registry lookup.
func (c *CheckHost) Whois(target string) (string, error) {
	req := TargetRequest{
		APIKey: c.APIKey,
		Target: target,
	}
	var whoisResponse string // Check-Host returns raw text/json loosely for whois.
	err := c.doRequest("POST", "whois", req, &whoisResponse)
	return whoisResponse, err
}

// Ping dispatches ICMP echo requests to the target from global nodes.
func (c *CheckHost) Ping(target string, options *MonitoringRequest) (*CheckCreated, error) {
	if options == nil {
		options = &MonitoringRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target

	var resp CheckCreated
	err := c.doRequest("POST", "ping", options, &resp)
	return &resp, err
}

// DNS queries global nameservers for specific DNS records.
func (c *CheckHost) DNS(target string, options *DNSTargetRequest) (*CheckCreated, error) {
	if options == nil {
		options = &DNSTargetRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target

	var resp CheckCreated
	err := c.doRequest("POST", "dns", options, &resp)
	return &resp, err
}

// TCP attempts to establish a 3-way TCP handshake on a specific destination port.
func (c *CheckHost) TCP(target string, port int, options *TCPMonitoringRequest) (*CheckCreated, error) {
	if options == nil {
		options = &TCPMonitoringRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target
	options.Port = port

	var resp CheckCreated
	err := c.doRequest("POST", "tcp", options, &resp)
	return &resp, err
}

// UDP sends UDP packets to a specified target and port.
func (c *CheckHost) UDP(target string, port int, options *UDPMonitoringRequest) (*CheckCreated, error) {
	if options == nil {
		options = &UDPMonitoringRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target
	options.Port = port

	var resp CheckCreated
	err := c.doRequest("POST", "udp", options, &resp)
	return &resp, err
}

// Http executes an HTTP/HTTPS request to the target to measure latency and statuses.
func (c *CheckHost) Http(target string, options *MonitoringRequest) (*CheckCreated, error) {
	if options == nil {
		options = &MonitoringRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target

	var resp CheckCreated
	err := c.doRequest("POST", "http", options, &resp)
	return &resp, err
}

// MTR initiates an MTR (My Traceroute) diagnostic.
func (c *CheckHost) MTR(target string, options *MTRMonitoringRequest) (*CheckCreated, error) {
	if options == nil {
		options = &MTRMonitoringRequest{}
	}
	options.APIKey = c.APIKey
	options.Target = target

	// Set default repeats based on swagger
	if options.RepeatChecks == 0 {
		options.RepeatChecks = 10
	}

	var resp CheckCreated
	err := c.doRequest("POST", "mtr", options, &resp)
	return &resp, err
}

// Report fetches the compiled report and real-time statuses from a previously initiated monitoring check via UUID.
func (c *CheckHost) Report(uuid string) (map[string]interface{}, error) {
	var report map[string]interface{}
	err := c.doRequest("GET", fmt.Sprintf("report/%s", uuid), nil, &report)
	return report, err
}
