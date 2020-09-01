package main

import (
	"flag"
	"github.com/go-kratos/kratos/pkg/net/trace/zipkin"
	appenv "github.com/vazmin/eagle-eye-kratos/common/env"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/api"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/di"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("licensing start")
	paladin.Init()
	api.Init()
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	zipkin.Init(&zipkin.Config{
		Endpoint: appenv.ZipkinEndpoint + "/api/v2/spans",
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("licensing exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
