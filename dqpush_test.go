package delay_queue

import (
	"fmt"
	"testing"

	"github.com/yasin-wu/delay-queue/delayqueue"
)

func TestDQPush(t *testing.T) {
	delayCmd := &delayqueue.DelayCmd{
		Config: &delayqueue.Config{
			BucketSize: 10,
			Redis: delayqueue.RedisConfig{
				Host:     "192.168.131.135:6379",
				PassWord: "1qazxsw21201",
			},
		},
	}
	delayCmd.Init()
	topic := "test_topic"
	for i := 1; i < 40; i++ {
		job := delayqueue.Job{
			Topic: topic,
			Id:    fmt.Sprintf("%d", i),
			Delay: int64(2 * i),
		}
		delayCmd.Push(job)
	}
}
