package main

import (
	"encoding/json"
	"log"
)

type Job struct {
	Topic   string     `json:"topic"`   //Job类型
	Id      string     `json:"id"`      //Job唯一标识
	Delay   int64      `json:"delay"`   //Job延迟时间
	TTR     int32      `json:"ttr"`     //Job延迟执行时间 /秒
	Body    string     `json:"body"`    //Job内容
	URL     string     `json:"url"`     //回调网址
}

//添加任务到Job Pool
//	key:   job.id
//  value: job
func putJobPool(key string, job *Job) error {
	value, err := json.Marshal(job)
	if err != nil {
		return err
	}
	//添加到Job Pool
	_, err = execRedisCommand("SET", key, value)
	if err != nil {
		log.Printf("SET命令执行失败_%s", key)
		return err
	}
	return nil
}

//根据job.id获取job
func getJobPool(key string) (*Job, error) {
	value, err := execRedisCommand("GET", key)
	if err != nil {
		log.Printf("GET命令执行失败_%s", key)
		return nil, err
	}

	if value == nil {
		return nil, nil
	}

	var job Job
	err = json.Unmarshal(value.([]byte), &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

//根据job.id删除存储的job元信息
func deleteJobPool(key string) error {
	_, err := execRedisCommand("DEL", key)
	return err
}