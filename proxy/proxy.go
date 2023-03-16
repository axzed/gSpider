package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

// roundRobinSwitcher 是一个轮询代理的结构体
type roundRobinSwitcher struct {
	ProxyURLs []*url.URL // 代理地址
	index     uint32     // 当前使用的代理地址的下标
}

// GetProxy 用于获取代理地址
func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	// 取余算法实现轮询调度
	// 每一次调用 GetProxy 函数，atomic.AddUint32 会将 index 加 1，并通过取余操作实现对代理地址的轮询
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.ProxyURLs[index%uint32(len(r.ProxyURLs))]
	return u, nil
}

// RoundRobinProxySwitcher creates a proxy switcher function which rotates
// ProxyURLs on every request.
// The proxy type is determined by the URL scheme. "http", "https"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func RoundRobinProxySwitcher(ProxyURLs ...string) (ProxyFunc, error) {
	// 如果代理地址列表为空，则返回错误
	if len(ProxyURLs) < 1 {
		return nil, errors.New("Proxy URL list is empty")
	}
	// 创建一个轮询代理的结构体
	urls := make([]*url.URL, len(ProxyURLs))
	// 遍历代理地址列表
	// 在遍历过程中，将代理地址解析为url.URL类型
	for i, u := range ProxyURLs {
		parseU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parseU
	}
	// 创建一个轮询代理的结构体
	r := &roundRobinSwitcher{urls, 0}
	// 返回一个代理函数
	proxyF := r.GetProxy
	return proxyF, nil
}
