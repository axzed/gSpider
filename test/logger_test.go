package test

import (
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func init() {
	// 设置日志前缀
	log.SetPrefix("TRACE: ")
	// 设置日志标志
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func TestLog(t *testing.T) {
	log.Println("message")
	log.Fatalln("fatal message")
	log.Panicln("panic message")
}

func TestZap(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "www.google.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
}

func TestSugZap(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	url := "www.google.com"
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
}
