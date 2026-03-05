package checkhost

type TargetRequest struct {
	APIKey string `json:"apikey,omitempty"`
	Target string `json:"target"`
}

type DNSTargetRequest struct {
	APIKey      string   `json:"apikey,omitempty"`
	Target      string   `json:"target"`
	QueryMethod string   `json:"querymethod,omitempty"`
	Region      []string `json:"region,omitempty"`
}

type MonitoringRequest struct {
	APIKey       string   `json:"apikey,omitempty"`
	Target       string   `json:"target"`
	RepeatChecks int      `json:"repeatchecks,omitempty"`
	Region       []string `json:"region,omitempty"`
	Timeout      int      `json:"timeout,omitempty"`
}

type TCPMonitoringRequest struct {
	MonitoringRequest
	Port int `json:"port"`
}

type UDPMonitoringRequest struct {
	MonitoringRequest
	Port    int    `json:"port"`
	Payload string `json:"payload,omitempty"`
}

type MTRMonitoringRequest struct {
	APIKey         string   `json:"apikey,omitempty"`
	Target         string   `json:"target"`
	RepeatChecks   int      `json:"repeatchecks,omitempty"`
	ForceIPVersion int      `json:"forceIPversion,omitempty"`
	ForceProtocol  string   `json:"forceProtocol,omitempty"`
	Region         []string `json:"region,omitempty"`
}

type MinResponseINFO struct {
	IP      string `json:"ip"`
	Reverse string `json:"reverse"`
	IPRange string `json:"iprange"`
	Country string `json:"country"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

type CheckCreated struct {
	Status       int    `json:"status"`
	Target       string `json:"target"`
	Method       string `json:"method"`
	RepeatChecks int    `json:"repeatchecks"`
	UUID         string `json:"uuid"`
	ReportURL    string `json:"reportURL"`
	APIURL       string `json:"apiURL"`
	Message      string `json:"message"`
	Success      bool   `json:"success"`
}

// Result mapping generic report data. Real world results return dicts mapped by Node name.
type Report struct {
	Nodes map[string]interface{}
}
