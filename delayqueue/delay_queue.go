package delayqueue

import (
	"github.com/yasin-wu/delay_queue/cronjob"
	"github.com/yasin-wu/delay_queue/logger"
	"github.com/yasin-wu/utils/redis"
)

var delayQueue *DelayQueue

type Option func(delayQueue *DelayQueue)

/**
 * @author: yasin
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
 * @author: yasin
 * @date: 2022/1/13 11:04
 * @params: host, keyPrefix string, batchLimit int, options ...redis.Option
 * @return: *DelayQueue
 * @description: 创建DelayQueue Client
 */
func New(host, keyPrefix string, batchLimit int, options ...redis.Option) *DelayQueue {
	delayQueue = &DelayQueue{}
	if keyPrefix == "" {
		keyPrefix = defaultKeyPrefix
	}
	if batchLimit == 0 {
		batchLimit = defaultBatchLimit
	}
	cli, err := redis.New(host, options...)
	if err != nil {
		return nil
	}
	redisCli := &redisClient{
		keyPrefix:  keyPrefix,
		batchLimit: batchLimit,
		client:     cli,
	}
	sche := cronjob.New()
	sche.Register([]int{}, 1, DelayQueueCronJob{})
	delayQueue.logger = logger.DefaultLogger
	delayQueue.scheduler = sche
	delayQueue.jobExecutorFactory = make(map[string]*jobExecutor)
	delayQueue.redisCli = redisCli
	return delayQueue
}

/**
 * @author: yasin
 * @date: 2022/1/13 11:05
 * @description: 启动延迟消息队列
 */
func (dq *DelayQueue) StartBackground() {
	dq.scheduler.Start()
}

/**
 * @author: yasin
 * @date: 2022/1/13 11:04
 * @params: action JobBaseAction
 * @return: error
 * @description: 注册延时服务
 */
func (dq *DelayQueue) Register(action JobBaseAction) error {
	jobID := action.ID()
	if _, ok := dq.jobExecutorFactory[jobID]; ok {
		return ErrorsDelayQueueRegisterIDDuplicate
	} else {
		dq.jobExecutorFactory[jobID] = &jobExecutor{
			ID:     jobID,
			action: action,
		}
	}
	return nil
}

/**
 * @author: yasin
 * @date: 2022/1/13 11:05
 * @params: job DelayJob
 * @return: error
 * @description: 添加延迟任务
 */
func (dq *DelayQueue) AddJob(job DelayJob) error {
	return dq.redisCli.zadd(job)
}

/**
 * @author: yasin
 * @date: 2022/1/13 13:22
 * @params: logger logger.Logger
 * @description: 设置日志
 */
func (dq *DelayQueue) SetLogger(logger logger.Logger) {
	dq.logger = logger
}

func (dq *DelayQueue) availableJobIDs() []string {
	var IDs []string
	for k := range dq.jobExecutorFactory {
		IDs = append(IDs, k)
	}
	return IDs
}
