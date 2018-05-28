package comm

import (
	"github.com/spf13/viper"
	"fmt"
	"upEletrcSign/config"
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

	if !config.HasModuleInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	pt := viper.GetInt("server.port")
	if pt <= 0 {
		return fmt.Errorf("port formate not crrect:%v", pt)
	}
	cf := &httpSvrConf{}
	cf.ListenIp = viper.GetString("server.host")
	cf.ListenPort = pt
	cf.RecvTimeOut = viper.GetInt("server.readTimeout")
	cf.WriteTimeOut = viper.GetInt("server.writeTimeout")
	t.conf = cf
	fmt.Println("Svr 加载成功")

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
