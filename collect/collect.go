package collect

import (
	"bufio"
	"fmt"
	"github.com/axzed/gSpider/proxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetcher is an interface that wraps the Get method.
type Fetcher interface {
	Get(url *Request) ([]byte, error)
}

// BaseFetch is a struct that implements Fetcher interface.
// BaseFetch is the basic logic of the fetch
type BaseFetch struct {
}

func (BaseFetch) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %v", resp.StatusCode)
		return nil, err
	}
	// 通过resp.Body创建一个bufio.Reader
	bodyReader := bufio.NewReader(resp.Body)
	// 用bodyReader创建一个utf-8的Reader
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	// 读取utf8Reader的内容
	return ioutil.ReadAll(utf8Reader)
}

// BrowserFetch 可以调整浏览器的header
type BrowserFetch struct {
	Timeout time.Duration   // 超时时间
	Proxy   proxy.ProxyFunc // 代理
}

// Get 用于获取网页的内容
func (b BrowserFetch) Get(request *Request) ([]byte, error) {

	// 新建一个http.Client,
	client := &http.Client{
		Timeout: b.Timeout, // 设置超时时间
	}
	if b.Proxy != nil {
		// 获取默认的transport
		transport := http.DefaultTransport.(*http.Transport)
		// 设置代理
		transport.Proxy = b.Proxy
		// 为客户端设置transport
		client.Transport = transport
	}
	// 新建一个http.Request
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
	}
	// 如果有cookie，就在请求头中设置cookie
	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}
	// 在请求头中设置header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 格式转换返回
	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// DeterminEncoding 用于判断网页的编码
func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("fetch error: %v", err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
