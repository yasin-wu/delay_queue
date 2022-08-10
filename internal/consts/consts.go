package consts

import (
	"errors"

	"github.com/yasin-wu/delay_queue/v2/pkg"
)

const (
	DefaultBatchLimit = 1000
	DefaultKeyPrefix  = "delay_queue"
)

var (
	ErrDelayQueueRegisterIDDuplicate = errors.New("your job id has been used")
	DefaultRedisOptions              = &pkg.Options{Addr: "localhost:6379", Password: "", DB: 0}
)
