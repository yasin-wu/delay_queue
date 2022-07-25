package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/yasin-wu/delay_queue/v2/delayqueue"
)

var redisOptions = &delayqueue.RedisOptions{Addr: "localhost:6379", Password: "yasinwu"}

type JobActionSMS struct{}

var _ delayqueue.JobBaseAction = (*JobActionSMS)(nil)

func (JobActionSMS) ID() string {
	return "JobActionSMS"
}

func (JobActionSMS) Execute(args []any) error {
	for _, arg := range args {
		if phoneNumber, ok := arg.(string); ok {
			fmt.Printf("sending sms to %s,time:%v\n", phoneNumber, time.Now().Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}

func TestDelayQueue(t *testing.T) {
	dq := delayqueue.New("test-yasin", 0, redisOptions)
	if err := dq.Register(JobActionSMS{}); err != nil {
		log.Fatal(err)
	}
	dq.StartBackground()
	fmt.Printf("add job:%v\n", time.Now().Format("2006-01-02 15:04:05"))
	if err := dq.AddJob(delayqueue.DelayJob{
		ID:        (&JobActionSMS{}).ID(),
		DelayTime: 5,
		Args:      []any{"181****9331"},
	}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Second)
}
