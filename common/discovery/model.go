package discovery

import "encoding/json"

// EndportInfo 机器信息：包含 IP、Port 和资源配置
type EndportInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	// MetaData 该机器的资源信息
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndportInfo, error) {
	ed := &EndportInfo{}
	err := json.Unmarshal(data, ed)
	if err != nil {
		return nil, err
	}
	return ed, nil
}

func (e *EndportInfo) Marshal() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(data)
}
