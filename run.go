package main

import (
	"net/http"
	"gopkg.in/ini.v1"
	"log"
)

var (
	port string   //web端口
)

func initRun() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("配置文件错误")
	}
	//设置运行模式
	runPattern = cfg.Section("run").Key("runPattern").String()
	//redis地址
	redisAddr = cfg.Section("redis").Key("addr").String()
	//web端口
	port = cfg.Section("run").Key("port").String()
	//delay bucket数量
	bucketNum, err = cfg.Section("run").Key("bucketNum").Int()
	if err != nil {
		log.Fatal("delay bucketNum设置错误")
	}
}

func about() {
	println(`
	 ____  _                   ____   ___  
	/ ___|| |_ ___  _ __   ___|  _ \ / _ \ 
	\___ \| __/ _ \| '_ \ / _ \ | | | | | |
	 ___) | || (_) | | | |  __/ |_| | |_| |
	|____/ \__\___/|_| |_|\___|____/ \__\_\

	Name: Stone Delay Queue		
	Author: StoneFlying	
	Email: stoneflying@yeah.net			   
	`)
}

func runWeb() {
	http.HandleFunc("/add", addHandle)
	http.HandleFunc("/pop", popHandle)
	http.HandleFunc("/finish", finishHandle)
	http.HandleFunc("/delete", deleteHandle)
	http.ListenAndServe(":"+port, nil)
}