package main

import (
	"strconv"
	"fmt"
	"log"
)

const (
	bucketName = "Bucket_%s"
)

var (
	bucketNameChan = make(chan string)
	bucketNum int
)

type jobItem struct {
	delay   int64     //到期时间
	id      string    //job.id
}

func initDelayBucket() {
	//通过轮询获取bucketName
	//将job id轮询放入bucket中
	go func() {
		index := 0
		for {
			bucketNameChan <- curBucket(index)
			index++
			index %= bucketNum
		}
	}()
}

//将job id放到redis的有序链表中
//根据socre进行排序
//此处score代表的是到期时间
func putBucket(key string, jobID string, score int64) error {
	_, err := execRedisCommand("ZADD", key, score, jobID)
	return err
}

//获取最近到期的job id
func getBucket(key string) (*jobItem, error) {
	value, err := execRedisCommand("ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		log.Printf("ZRANGE命令执行失败_%s", key)
		return nil, err
	}

	if len(value.([]interface{})) != 2 {
		return nil, nil
	}

	delayStr := string((value.([]interface{}))[1].([]byte))
	delay, _ := strconv.ParseInt(delayStr, 10, 64)
	item := &jobItem{
		delay: delay,
		id:    string((value.([]interface{}))[0].([]byte)),
	}

	return item, nil
}

//删除delay bucket中的指定job
func deleteBucket(key, jobID string) error {
	_, err := execRedisCommand("ZREM", key, jobID)
	if err != nil {
		log.Printf("ZREM命令执行失败_%s", key)
		return err
	}
	return nil
}

func curBucket(index int) string {
	return fmt.Sprintf(bucketName, strconv.Itoa(index))
}