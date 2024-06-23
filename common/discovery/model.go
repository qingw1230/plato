package discovery

import "encoding/json"

// EndpointInfo 机器信息：包含 IP、Port 和资源配置
type EndpointInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	// MetaData 该机器的资源信息
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	ed := &EndpointInfo{}
	err := json.Unmarshal(data, ed)
	if err != nil {
		return nil, err
	}
	return ed, nil
}

func (e *EndpointInfo) Marshal() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(data)
}
