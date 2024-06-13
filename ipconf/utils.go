package ipconf

import "github.com/qingw1230/plato/ipconf/domain"

// top5Endports 获取得分前 5 的机器
func top5Endports(eds []*domain.Endport) []*domain.Endport {
	if len(eds) < 5 {
		return eds
	}
	return eds[:5]
}

// packRes 将机器列表信息包装成响应体
func packRes(ed []*domain.Endport) Response {
	return Response{
		Message: "ok",
		Code:    0,
		Data:    ed,
	}
}
