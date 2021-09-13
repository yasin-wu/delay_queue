package delayqueue

import (
	"errors"
)

const (
	defaultBatchLimit = 1000
	defaultKeyPrefix  = "delay_queue"
)

var (
	ErrorsDelayQueueRegisterIDDuplicate = errors.New("your job id has been used")
)
