package delay_queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/yasin-wu/delay-queue/delayqueue"
)

func TestDQPop(t *testing.T) {
	delayCmd := &delayqueue.DelayCmd{
		Config: &delayqueue.Config{
			BucketSize: 10,
			Redis: delayqueue.RedisConfig{
				Host:     "192.168.131.135:6379",
				Password: "1qazxsw21201",
			},
		},
	}
	delayCmd.Init()
	topic := "test_topic"
	for {
		job, err := delayCmd.Pop(topic)
		if err != nil {
			t.Error(err.Error())
		}
		if job == nil {
			continue
		}
		fmt.Println(fmt.Sprintf("id:%s,time:%v", job.Id, time.Now()))
		delayCmd.Delete(job.Id)
		fmt.Println(fmt.Sprintf("del:%s", job.Id))
	}
}
