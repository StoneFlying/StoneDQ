package main

import (
	"log"
	"time"
	"encoding/json"
)

var (
	timers []*time.Ticker
)

func initTimer() {
	timers = make([]*time.Ticker, bucketNum)
	for k := range timers {
		bucketName := curBucket(k)
		timers[k] = time.NewTicker(1 * time.Second)
		go ticker(timers[k], bucketName)
	}
}

func ticker(t *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-t.C:
			findReadyJobFromBucket(t, bucketName)
		}
	}
}

func findReadyJobFromBucket(t time.Time, bucketName string) {
	//从job bucket中取出头job
	jobItem, err := getBucket(bucketName)
	if err != nil {
		log.Println(err)
		return
	}

	if jobItem == nil {
		return
	}

	//println(bucketName, jobItem.delay, jobItem.id)

	//延迟时间还未到达
	if jobItem.delay > t.Unix() {
		return
	}

	//从job pool中取出元信息
	job, err := getJobPool(jobItem.id)
	if err != nil {
		log.Println(err)
		return
	}

	//job处于delete状态，直接pass
	//并从delay bucket中进行删除
	if job == nil {
		deleteBucket(bucketName, jobItem.id)
		return
	}

	//再次确认delay
	//以防delay bucket中或者job pool中信息被篡改后导致的数据不一致
	if job.Delay >= t.Unix() {
		//删除delay bucket中存储的该任务
		deleteBucket(bucketName, job.Id)
		//根据job.Delay中的延时信息重新放入delay bucket
		putBucket(<-bucketNameChan, job.Id, job.Delay)
		return
	}

	//轮询处理
	if runPattern == "polling" {
		//放入ready queue，等待客户端消费
		putReadyQueue(job.Topic, job.Id)
		//从delay bucket中移除
		deleteBucket(bucketName, job.Id)
	} else if runPattern == "callback" {  //回调处理
		value, err := json.Marshal(job)
		if err != nil {
			return
		}
		str, err := httpPost(job.URL, "data="+string(value))
		if err != nil {
			return
		}
		if str == "success" {
			//从delay bucket中移除
			deleteBucket(bucketName, job.Id)
			//从job pool移除
			deleteJobPool(job.Id)
			log.Printf("%s_【执行成功】_返回内容:%s", string(value), str)
		} else {
			//修改延时时间后放入bucket，避免失败回调一直重试
			//默认延迟10秒后重试
			deleteBucket(bucketName, job.Id)
			putBucket(<-bucketNameChan, job.Id, t.Unix() + 10)
			log.Printf("%s___【执行失败】___返回内容:%s", string(value), str)
		}
	}
}