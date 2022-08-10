package dqueue

import (
	"github.com/yasin-wu/delay_queue/v2/internal/dqueue"
	"github.com/yasin-wu/delay_queue/v2/internal/redis"

	"github.com/yasin-wu/delay_queue/v2/pkg"

	"github.com/yasin-wu/delay_queue/v2/internal/consts"

	"github.com/yasin-wu/delay_queue/v2/internal/logger"

	"github.com/yasin-wu/delay_queue/v2/internal/cronjob"
)

var delayQueue *DelayQueue

type DelayQueue struct {
	logger             logger.Logger
	scheduler          *cronjob.Scheduler
	jobExecutorFactory map[string]*dqueue.JobExecutor
	redisCli           *redis.Client
}

func New(keyPrefix string, batchLimit int64, options *pkg.Options) *DelayQueue {
	if options == nil {
		options = consts.DefaultRedisOptions
	}
	if keyPrefix == "" {
		keyPrefix = consts.DefaultKeyPrefix
	}
	if batchLimit == 0 {
		batchLimit = consts.DefaultBatchLimit
	}
	delayQueue = &DelayQueue{}
	sche := cronjob.New()
	sche.Register([]int{}, 1, CronJob{})
	delayQueue.logger = logger.DefaultLogger
	delayQueue.scheduler = sche
	delayQueue.jobExecutorFactory = make(map[string]*dqueue.JobExecutor)
	delayQueue.redisCli = redis.New(keyPrefix, batchLimit, options)
	return delayQueue
}

func (dq *DelayQueue) StartBackground() {
	dq.scheduler.Start()
}

func (dq *DelayQueue) Register(action pkg.JobBaseAction) error {
	jobID := action.ID()
	if _, ok := dq.jobExecutorFactory[jobID]; ok {
		return consts.ErrDelayQueueRegisterIDDuplicate
	}
	dq.jobExecutorFactory[jobID] = &dqueue.JobExecutor{
		ID:     jobID,
		Action: action,
	}
	return nil
}

func (dq *DelayQueue) AddJob(job pkg.DelayJob) error {
	return dq.redisCli.Zadd(job)
}

func (dq *DelayQueue) SetLogger(logger logger.Logger) {
	dq.logger = logger
	dq.redisCli.SetLogger(logger)
}

func (dq *DelayQueue) availableJobIDs() []string {
	var ids []string
	for k := range dq.jobExecutorFactory {
		ids = append(ids, k)
	}
	return ids
}
