package main

import (
	"runtime"
	"upEletrcSign/config"
	"fmt"
	"os"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"upEletrcSign/handle"
	"upEletrcSign/comm"
	"upEletrcSign/logr"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化配置文件
	err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//初始化配置文件
	err = logr.InitModules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/upSign", handle.DoHandle)
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(r)
	sockSvr := comm.HttpSvr{}
	err = sockSvr.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sockSvr.RunSvr(n)

	os.Exit(0)
}
