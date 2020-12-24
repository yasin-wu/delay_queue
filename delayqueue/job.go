package delayqueue

import (
	"encoding/json"
)

type Job struct {
	Topic string `json:"topic"`
	Id    string `json:"id"`    // job唯一标识ID
	Delay int64  `json:"delay"` // 延迟时间, unix时间戳
	TTR   int64  `json:"ttr"`
	Body  string `json:"body"` // job 对应内容
}

func (this *Job) getJob(key string) (*Job, error) {
	value, err := execRedisCommand("GET", key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, ErrRecordNotFound
	}
	byteValue := value.([]byte)
	job := &Job{}
	err = json.Unmarshal(byteValue, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (this *Job) putJob(key string, job Job) error {
	value, err := json.Marshal(job)
	if err != nil {
		return err
	}
	_, err = execRedisCommand("SET", key, value, "EX", Setting.Redis.ExpireTime)
	return err
}

func (this *Job) removeJob(key string) error {
	_, err := execRedisCommand("DEL", key)
	return err
}
