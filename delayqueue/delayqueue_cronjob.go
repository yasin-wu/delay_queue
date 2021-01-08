package delayqueue

type DelayQueueCronJob struct{}

func (DelayQueueCronJob) Name() string {
	return "DelayQCron"
}

func (DelayQueueCronJob) Process() error {
	IDs := delayQueue.availableJobIDs()
	return delayQueue.redisCli.BatchHandle(IDs)
}

func (DelayQueueCronJob) IfActive() bool {
	return true
}

func (DelayQueueCronJob) IfReboot() bool {
	return true
}
