package discovery

import "encoding/json"

// Endpoint 机器信息：包含 IP、Port 和资源信息
type EndpointInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	// MetaData 该机器的资源信息
	MetaData map[string]interface{} `json:"meta"`
}

func (e *EndpointInfo) Marshal() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	e := &EndpointInfo{}
	err := json.Unmarshal(data, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
