package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type response struct {
	Success    bool   `json:"success"`
	Error      string  `json:"err"`   //error reson
	ID         string `json:"id"`
	Value      string `json:"value"`   //job body
}

// 添加延迟任务
func addHandle(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		return
	}

	var job Job
	err = json.Unmarshal(body, &job)
	if err != nil {
		w.Write(newFailedResponse("", err.Error()))
		return
	}

	if runPattern == "callback" {
		if strings.TrimSpace(job.URL) == "" {
			w.Write(newFailedResponse(job.Id, "callback模式下回调url不能为空"))
			return
		}
	}

	if strings.TrimSpace(job.Topic) == "" || strings.TrimSpace(job.Body) == "" || 
		strings.TrimSpace(job.Id) == "" {
		w.Write(newFailedResponse(job.Id, "topic/body/id/均不能留空"))
		return
	}

	if job.Delay <= 0 || job.TTR <= 0 {
		w.Write(newFailedResponse(job.Id, "delay/ttr均不能小于0"))
		return
	}

	err = add(&job)
	if (err != nil) {
		w.Write(newFailedResponse(job.Id, err.Error()))
		return
	}
	w.Write(newSuccessResponse(job.Id, "添加成功"))
}

// 获取延迟任务
func popHandle(w http.ResponseWriter, r *http.Request) {
	if runPattern != "polling" {
		w.Write(newFailedResponse("", "非polling模式下，不能主动获取任务"))
		return
	}
	body, err := readBody(r)
	if err != nil {
		return
	}

	topic := &struct {
		Topic string `json:"topic"`
	}{
		Topic: "",
	}
	err = json.Unmarshal(body, topic)
	if err != nil {
		w.Write(newFailedResponse("", err.Error()))
		return
	}

	if (strings.TrimSpace(topic.Topic) == "") {
		w.Write(newFailedResponse("", "topic不能为空"))
		return
	}

	job, err := pop(topic.Topic)
	if err != nil {
		w.Write(newFailedResponse("", err.Error()))
		return
	}
	if job == nil {
		w.Write(newFailedResponse("", "暂无就绪任务"))
		return
	}
	w.Write(newSuccessResponse(job.Id, job.Body))
}

// 删除未执行延迟任务
func deleteHandle(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		return
	}

	id := &struct {
		Id string `json:"id"`
	}{
		Id: "",
	}
	err = json.Unmarshal(body, id)
	if err != nil {
		w.Write(newFailedResponse("", err.Error()))
		return
	}

	if (strings.TrimSpace(id.Id) == "") {
		w.Write(newFailedResponse("", "id不能为空"))
		return
	}

	err = delete(id.Id)
	if err != nil {
		w.Write(newFailedResponse("", err.Error()))
		return
	}
	w.Write(newSuccessResponse(id.Id, "删除成功"))
}

func finishHandle(w http.ResponseWriter, r *http.Request) {
	deleteHandle(w, r)
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func newSuccessResponse(id string, value string) []byte {
	return newResponse(true, "", id, value)
}

func newFailedResponse(id string, error string) []byte {
	return newResponse(false, error, id, "")
}

func newResponse(success bool, error string, id string, value string) []byte {
	resp := &response {
		Success: success, 
		Error:   error, 
		ID:      id, 
		Value:   value,
	}
	jsIndent,_ := json.MarshalIndent(resp, "", "\t")
	return jsIndent
}

// 发起post请求，服务器通过job数据主动回调相应url
// 回调url需要返回success来代表操作成功
func httpPost(url, data string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//if (string(body) != "success") {
	//	return "", nil
	//}
	return string(body), nil
}