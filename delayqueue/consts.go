package delayqueue

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	defaultBatchLimit = 1000
	defaultKeyPrefix  = "delay_queue"
)

const (
	defaultRedisHost           = "127.0.0.1:6379"
	defaultRedisPassWord       = ""
	defaultRedisDB             = 0
	defaultRedisMaxIdle        = 10
	defaultRedisMaxActive      = 0
	defaultRedisConnectTimeout = 5000
	defaultRedisReadTimeout    = 180000
	defaultRedisWriteTimeout   = 3000
	defaultIdleTimeout         = 300 * time.Second
)

var (
	redisCli   *redisClient
	redisPool  *redis.Pool
	delayQueue *DelayQueue
)

var (
	ErrorsDelayQueueRegisterIDDuplicate = errors.New("your job id has been used")
)
