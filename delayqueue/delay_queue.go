package delayqueue

import (
	"errors"
	"time"
)

type DelayQueue struct {
	job        Job
	bucket     BucketItem
	readyQueue ReadyQueue
}

func (this *DelayQueue) push(job Job) error {
	if job.Id == "" || job.Topic == "" || job.Delay < 0 {
		return errors.New("invalid job")
	}

	err := this.job.putJob(job.Id, job)
	if err != nil {
		return err
	}
	err = this.bucket.pushToBucket(<-bucketNameChan, job.Delay, job.Id)
	if err != nil {
		return err
	}

	return nil
}

func (this *DelayQueue) pop(topics []string) (*Job, error) {
	jobId, err := this.readyQueue.blockPopFromReadyQueue(topics, Setting.QueueBlockTimeout)
	if err != nil {
		return nil, err
	}

	if jobId == "" {
		return nil, errors.New("queue is nil")
	}

	job, err := this.job.getJob(jobId)
	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, errors.New("job is nil")
	}

	//timestamp := time.Now().Unix() + job.Delay
	//err = this.bucket.pushToBucket(<-bucketNameChan, timestamp, job.Id)
	return job, err
}

func (this *DelayQueue) remove(jobId string) error {
	return this.job.removeJob(jobId)
}

func (this *DelayQueue) get(jobId string) (*Job, error) {
	return this.job.getJob(jobId)
}

func waitTicker(timer *time.Ticker, bucketName string) {
	this := &DelayQueue{}
	for {
		select {
		case t := <-timer.C:
			this.tickHandler(t, bucketName)
		}
	}
}

func (this *DelayQueue) tickHandler(t time.Time, bucketName string) {
	for {
		bucketItem, err := this.bucket.getFromBucket(bucketName)
		if err != nil {
			return
		}

		if bucketItem == nil {
			return
		}

		// 延迟时间未到
		if bucketItem.timestamp > t.Unix() {
			return
		}

		// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		job, err := this.job.getJob(bucketItem.jobId)
		if err != nil {
			continue
		}

		// job元信息不存在, 从bucket中删除
		if job == nil {
			this.bucket.removeFromBucket(bucketName, bucketItem.jobId)
			continue
		}

		// 再次确认元信息中delay是否小于等于当前时间
		if job.Delay > t.Unix() {
			// 从bucket中删除旧的jobId
			this.bucket.removeFromBucket(bucketName, bucketItem.jobId)
			// 重新计算delay时间并放入bucket中
			this.bucket.pushToBucket(<-bucketNameChan, job.Delay, bucketItem.jobId)
			continue
		}

		err = this.readyQueue.pushToReadyQueue(job.Topic, bucketItem.jobId)
		if err != nil {
			continue
		}

		// 从bucket中删除
		this.bucket.removeFromBucket(bucketName, bucketItem.jobId)
	}
}
