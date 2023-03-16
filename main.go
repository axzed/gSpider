package main

import (
	"github.com/axzed/gSpider/collect"
	"github.com/axzed/gSpider/log"
	"github.com/axzed/gSpider/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regexp"
	"time"
)

var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}
	url := "https://google.com"
	var f collect.Fetcher = collect.BrowserFetch{
		TimeOut: 3000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		logger.Error("read content failed",
			zap.Error(err),
		)
		return
	}
	logger.Info("get content", zap.Int("len", len(body)))
}
