package collect

// Request 中包含了一个 URL，表示要访问的网站。
// 这里的 ParseFunc 函数会解析从网站获取到的网站信息
// 返回 Requesrts 数组用于进一步获取数据。
type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte, *Request) ParseResult
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}
