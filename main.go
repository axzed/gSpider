package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"time"
)

var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	// 1. 创建谷歌浏览器实例
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// 2. 设置context超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 3. 爬取网页内容
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.zhihu.com/question/381784377`),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\\n%s", example)
}
