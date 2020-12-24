package delayqueue

import (
	"fmt"
)

type ReadyQueue struct {
}

func (this *ReadyQueue) pushToReadyQueue(queueName string, jobId string) error {
	queueName = fmt.Sprintf(Setting.QueueName, queueName)
	_, err := execRedisCommand("RPUSH", queueName, jobId)
	return err
}

// 从队列中阻塞获取JobId
func (this *ReadyQueue) blockPopFromReadyQueue(queues []string, timeout int) (string, error) {
	var args []interface{}
	for _, queue := range queues {
		queue = fmt.Sprintf(Setting.QueueName, queue)
		args = append(args, queue)
	}
	args = append(args, timeout)
	value, err := execRedisCommand("BLPOP", args...)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", ErrRecordNotFound
	}
	var valueBytes []interface{}
	valueBytes = value.([]interface{})
	if len(valueBytes) == 0 {
		return "", ErrRecordNotFound
	}
	element := string(valueBytes[1].([]byte))
	return element, nil
}
