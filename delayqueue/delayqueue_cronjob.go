package delayqueue

type CronJob struct{}

func (CronJob) Name() string {
	return "DelayQCron"
}

func (CronJob) Process() error {
	IDs := delayQueue.availableJobIDs()
	return delayQueue.redisCli.batchHandle(IDs)
}

func (CronJob) IfActive() bool {
	return true
}

func (CronJob) IfReboot() bool {
	return true
}
