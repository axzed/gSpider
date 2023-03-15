package collect

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

// Fetcher is an interface that wraps the Get method.
type Fetcher interface {
	Get(url string) ([]byte, error)
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
}

// Get 用于获取网页的内容
func (BrowserFetch) Get(url string) ([]byte, error) {
	// 新建一个http.Client,
	client := &http.Client{}
	// 新建一个http.Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
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
