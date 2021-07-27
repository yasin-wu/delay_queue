package delayqueue

import "github.com/yasin-wu/utils/redis"

type Config struct {
	KeyPrefix  string
	BatchLimit int64
	Redis      *redis.Config
}
