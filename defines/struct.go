package defines

type Jwt struct {
}

type DnsPacket struct {
}

// 开发时采用的标准JSON格式的HTTP响应
type JsonHttpRes struct {
	Code int    // 响应码
	Msg  string // 信息
	Data any    // 数据
}
