package delayqueue

import (
	"github.com/yasin-wu/delay-queue/cronjob"
	"github.com/yasin-wu/delay-queue/logger"
	"github.com/yasin-wu/utils/redis"
)

var (
	redisCli   *redisClient
	delayQueue *DelayQueue
)

func initDelayQueue(conf *Config) {
	delayQueue = new(DelayQueue)
	sche := cronjob.New()
	sche.Register([]int{}, 1, DelayQueueCronJob{})
	delayQueue.logger = logger.DefaultLogger
	delayQueue.scheduler = sche
	delayQueue.jobExecutorFactory = make(map[string]*jobExecutor)
	delayQueue.redisCli = initRedis(conf)
	delayQueue.redisConf = conf.Redis
	delayQueue.logger.InfoF("Initialization of DelayQueue completed......")
}

func initRedis(conf *Config) *redisClient {
	if conf.KeyPrefix == "" {
		conf.KeyPrefix = defaultKeyPrefix
	}
	if conf.BatchLimit == 0 {
		conf.BatchLimit = defaultBatchLimit
	}
	client, err := redis.New(conf.Redis)
	if err != nil {
		panic(err)
	}
	redisCli = &redisClient{
		keyPrefix:  conf.KeyPrefix,
		batchLimit: conf.BatchLimit,
		client:     client,
	}
	return redisCli
}
