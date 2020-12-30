package delayqueue

import (
	"errors"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type DelayCmd struct {
	Config     *Config
	delayQueue DelayQueue
}

// Push 添加job
func (this *DelayCmd) Push(job Job) error {
	job.Id = strings.TrimSpace(job.Id)
	job.Topic = strings.TrimSpace(job.Topic)
	job.Body = strings.TrimSpace(job.Body)
	if job.Id == "" {
		return errors.New("job id不能为空")
	}
	if job.Topic == "" {
		return errors.New("topic 不能为空")
	}

	if job.Delay <= 0 || job.Delay > (1<<31) {
		return errors.New("delay 取值范围1 - (2^31 - 1)")
	}

	if job.TTR <= 0 || job.TTR > 86400 {
		return errors.New("ttr 取值范围1 - 86400")
	}

	job.Delay = time.Now().Unix() + job.Delay
	err := this.delayQueue.push(job)

	if err != nil {
		return err
	}
	return nil
}

// Pop 获取job
func (this *DelayCmd) Pop(topic string) (*Job, error) {

	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, errors.New("topic不能为空")
	}
	topics := strings.Split(topic, ",")
	job, err := this.delayQueue.pop(topics)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job is nil")
	}

	return job, nil
}

// Delete 删除job
func (this *DelayCmd) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("job id不能为空")
	}

	err := this.delayQueue.remove(id)
	if err != nil {
		return err
	}
	return nil
}

// Get 查询job
func (this *DelayCmd) Get(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id不能为空")
	}
	job, err := this.delayQueue.get(id)
	if err != nil {
		return err
	}

	if job == nil {
		return nil
	}
	return nil
}

func (this *DelayCmd) ClearCache() error {
	res, err := redis.Strings(execRedisCommand("KEYS", "*"))
	if err != nil {
		return err
	}
	for _, v := range res {
		_, err := execRedisCommand("DEL", v)
		if err != nil {
			continue
		}
	}
	return nil
}
