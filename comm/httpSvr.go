package comm

import (
	"fmt"
	"mygolib/modules/config"
	"mygolib/modules/myLogger"
	"errors"
	"time"
	"os"
	"os/signal"
	"net/http"
	"log"
	"context"
)

type httpSvrConf struct {
	ListenIp     string
	ListenPort   int
	RecvTimeOut  int
	WriteTimeOut int
	MaxAccNum    int
}

type HttpSvr struct {
	conf *httpSvrConf
}

func (t *HttpSvr) InitConfig() error {

	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}
	if !myLogger.HasLoggerInit() {
		return errors.New("日志模块未初始化，请先初始化")
	}

	config.SetSection("server")

	cf := &httpSvrConf{}
	cf.ListenIp = config.StringDefault("host", "")
	cf.ListenPort = config.IntDefault("port", 9090)
	cf.RecvTimeOut = config.IntDefault("readTimeout", 30)
	cf.WriteTimeOut = config.IntDefault("writeTimeout", 30)
	t.conf = cf
	fmt.Println("HttpSvr加载成功")

	return nil
}
func (t *HttpSvr) RunSvr(h http.Handler) {
	srv := &http.Server{
		Addr:           t.conf.ListenIp + fmt.Sprintf(":%d", t.conf.ListenPort),
		Handler:        h,
		ReadTimeout:    time.Duration(t.conf.RecvTimeOut) * time.Second,
		WriteTimeout:   time.Duration(t.conf.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		defer myLogger.Info("----HttpSvr关闭----")
		myLogger.Info("----HttpSvr启动----")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
}
