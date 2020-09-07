# StoneDQ
Go+Redis构建的延迟队列，包括主动拉取和回调两种模式
参考实现：[有赞延迟队列设计](https://tech.youzan.com/queuing_delay/)

         ____  _                   ____   ___  
        / ___|| |_ ___  _ __   ___|  _ \ / _ \ 
        \___ \| __/ _ \| '_ \ / _ \ | | | | | |
         ___) | || (_) | | | |  __/ |_| | |_| |
        |____/ \__\___/|_| |_|\___|____/ \__\_\

        Name: Stone Delay Queue                
        Author: StoneFlying
        Email: stoneflying@yeah.net
        
## polling模式下
# 使用
### /add  
#### request
```
{
	"topic": "订单",
	"id": "1",
	"delay": 10,
	"ttr": 10,
	"body": "自定义请求内容"
}
```
#### response
```
{
	"success": true,
	"err": "",
	"id": "1",
	"value": "添加成功"
}
```

### /pop
#### request
```
{
	"topic": "订单"
}
```
#### response
```
{
	"success": true,
	"err": "",
	"id": "1",
	"value": "自定义请求内容"
}
```

### /finish
#### request
```
{
	"id": "1"
}
```
#### response
```
{
	"success": true,
	"err": "",
	"id": "1",
	"value": "删除成功"
}
```

### /delete
#### request
```
{
	"id": "1"
}
```
#### response
```
{
	"success": true,
	"err": "",
	"id": "1",
	"value": "删除成功"
}
```

## callback模式下
### /add  
#### request
```
{
	"topic": "订单",
	"id": "1",
	"delay": 10,
	"ttr": 10,
	"body": "自定义请求内容",
	"url": "http://localhost"
}
```
#### response
```
{
	"success": true,
	"err": "",
	"id": "1",
	"value": "添加成功"
}
```
在callback模式下服务器会自动回调request中"url": "http://localhost"指定的地址  
本例回调时传递的POST参数:{"topic":"订单","id":"1","delay":1599493747,"ttr":10,"body":"自定义请求内容","url":"http://localhost"}  
客户端执行成功应输出:success，以告知客户端执行成功，以便服务器将数据从延迟队列删除  

## 注意：
Job id需要在所有业务组内全局唯一，建议通过topic+id构成  
在polling模式下，客户端获取数据处理完成后，需要访问/finish告知服务器客户机处理成功  
服务端保证数据至少消费一次，在只能消费一次的情况下，需要客户端自行保证  
