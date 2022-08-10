package dqueue

import (
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/yasin-wu/delay_queue/v2/pkg"
)

type CronJob struct{}

func (CronJob) Name() string {
	return "DelayQCron"
}

func (CronJob) Process() error {
	return delayQueue.batchHandle(delayQueue.availableJobIDs())
}

func (CronJob) IfActive() bool {
	return true
}

func (CronJob) IfReboot() bool {
	return true
}

func (d *DelayQueue) batchHandle(ids []string) error {
	var wg = sync.WaitGroup{}
	wg.Add(len(ids))
	for _, name := range ids {
		key := d.redisCli.FormatKey(name)
		go func(key string) {
			batch, lastScore, err := d.redisCli.GetBatch(key)
			if err != nil {
				d.logger.Errorf("get batch failed , error:%s", err.Error())
				return
			}
			d.executeBatch(batch)
			defer func() {
				if err != nil || len(batch) != 0 {
					d.redisCli.ClearBatch(key, lastScore)
				}
				wg.Done()
			}()
		}(key)
	}
	wg.Wait()
	return nil
}

func (d *DelayQueue) executeBatch(batch []redis.Z) {
	for _, item := range batch {
		var delayJob pkg.DelayJob
		if err := json.Unmarshal([]byte(item.Member.(string)), &delayJob); err != nil {
			d.logger.Errorf("json unmarshal failed , error:%s", err.Error())
			continue
		}
		executor, ok := d.jobExecutorFactory[delayJob.ID]
		if !ok {
			continue
		}
		if err := executor.Action.Execute(delayJob.Args); err != nil {
			d.logger.Errorf("job action execute failed , error:%s", err.Error())
		}
	}
}
