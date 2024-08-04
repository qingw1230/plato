package discov

// Service 服务和提供该服务的机器列表
type Service struct {
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

// Endpoint 机器的相关信息
type Endpoint struct {
	ServerName string `json:"server_name"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
	Enable     bool   `json:"enable"`
}
