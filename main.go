package main

var (
	runPattern string  //运行模式
)

func main() {
	about()
	//初始化站点配置
	initRun()
	//初始化redis
	initRedis()
	//初始化delay bucket
	initDelayBucket()
	//初始化timer
	initTimer()
	//运行http服务
	runWeb()
}