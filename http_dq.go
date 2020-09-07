package main

import (
	"time"
)

//const (
//	jobPoolName = "JOB_POOL"
//)

//put待延迟任务
func add(job *Job) error {
	//计算绝对执行时间
	job.Delay += time.Now().Unix()

	//放入job pool
	err := putJobPool(job.Id, job)
	if err != nil {
		return err
	}

	//放入delay bucket
	err = putBucket(<-bucketNameChan, job.Id, job.Delay)
	if err != nil {
		return err
	}
	return nil
}

//取出以到达延迟时间任务返回给客户端
func pop(jobTopic string) (*Job, error) {
	//从ready queue获取已经ready的job
	jobID, err := getReadyQueue(jobTopic)
	if err != nil {
		return nil, err
	}

	if jobID == "" {
		return nil, nil
	}

	//根据job id 查询job元信息
	job, err := getJobPool(jobID)
	if err != nil {
		return nil, err
	}

	//job已被删除
	if job == nil {
		return nil, nil
	}

	//重新计算执行时间，放入delay bucket
	//确保每个job至少被消费一次
	delay := time.Now().Unix() + int64(job.TTR)
	job.Delay = delay
	putBucket(<-bucketNameChan, job.Id, delay)

	return job, nil
}

//客户端操作完成，通过finish通知服务端
//服务端删除job对应元信息
func finish(jobID string) error {
	return deleteJobPool(jobID)
}

//客户端直接删除对应job元信息
func delete(jobID string) error {
	return deleteJobPool(jobID)
}