//延迟时间已经到达的任务
package main

import (
	"log"
)

//已到达时间的延迟任务放入相应topic的ReadyQueue
func putReadyQueue(key, jobID string) error {
	_, err := execRedisCommand("RPUSH", key, jobID)
	if err != nil {
		log.Printf("RPUSH命令执行失败_%s", key)
		return err
	}
	
	return nil
}

//取出指定topic中的已到达延迟任务
func getReadyQueue(key string) (string, error) {
	jobID, err := execRedisCommand("LPOP", key)
	if err != nil {
		log.Printf("LPOP命令执行失败_%s", key)
		return "", err
	}
	if jobID == nil {
		return "", nil
	}

	return string(jobID.([]byte)), nil
}