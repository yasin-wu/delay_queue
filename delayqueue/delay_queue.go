package delayqueue

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/yasin-wu/delay_queue/v2/cronjob"
	"github.com/yasin-wu/delay_queue/v2/logger"
)

var delayQueue *DelayQueue

type Option func(delayQueue *DelayQueue)

/**
 * @author: yasinWu
 * @date: 2022/1/13 11:03
 * @description: DelayQueue Client
 */
type DelayQueue struct {
	logger             logger.Logger
	scheduler          *cronjob.Scheduler
	jobExecutorFactory map[string]*jobExecutor
	redisCli           *redisClient
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 11:04
 * @params: keyPrefix string, batchLimit int64, redisOptions *RedisOptions
 * @return: *DelayQueue
 * @description: 创建DelayQueue Client
 */
func New(keyPrefix string, batchLimit int64, redisOptions *RedisOptions) *DelayQueue {
	if redisOptions == nil {
		redisOptions = defaultRedisOptions
	}
	delayQueue = &DelayQueue{}
	if keyPrefix == "" {
		keyPrefix = defaultKeyPrefix
	}
	if batchLimit == 0 {
		batchLimit = defaultBatchLimit
	}
	redisCli := &redisClient{
		keyPrefix:  keyPrefix,
		batchLimit: batchLimit,
		client:     redis.NewClient((*redis.Options)(redisOptions)),
		ctx:        context.Background(),
	}
	sche := cronjob.New()
	sche.Register([]int{}, 1, CronJob{})
	delayQueue.logger = logger.DefaultLogger
	delayQueue.scheduler = sche
	delayQueue.jobExecutorFactory = make(map[string]*jobExecutor)
	delayQueue.redisCli = redisCli
	return delayQueue
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 11:05
 * @description: 启动延迟消息队列
 */
func (dq *DelayQueue) StartBackground() {
	dq.scheduler.Start()
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 11:04
 * @params: action JobBaseAction
 * @return: error
 * @description: 注册延时服务
 */
func (dq *DelayQueue) Register(action JobBaseAction) error {
	jobID := action.ID()
	_, ok := dq.jobExecutorFactory[jobID]
	if ok {
		return ErrorsDelayQueueRegisterIDDuplicate
	}
	dq.jobExecutorFactory[jobID] = &jobExecutor{
		ID:     jobID,
		action: action,
	}
	return nil
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 11:05
 * @params: job DelayJob
 * @return: error
 * @description: 添加延迟任务
 */
func (dq *DelayQueue) AddJob(job DelayJob) error {
	return dq.redisCli.zadd(job)
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 13:22
 * @params: logger logger.Logger
 * @description: 设置日志
 */
func (dq *DelayQueue) SetLogger(logger logger.Logger) {
	dq.logger = logger
}

func (dq *DelayQueue) availableJobIDs() []string {
	var IDs []string //nolint:prealloc
	for k := range dq.jobExecutorFactory {
		IDs = append(IDs, k)
	}
	return IDs
}
