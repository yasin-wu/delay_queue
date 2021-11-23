package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	"github.com/yasin-wu/utils/redis"

	"github.com/yasin-wu/delay_queue/delayqueue"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type JobActionSMS struct{}

var _ delayqueue.JobBaseAction = (*JobActionSMS)(nil)

func (JobActionSMS) ID() string {
	return "JobActionSMS"
}

func (JobActionSMS) Execute(args []interface{}) error {
	for _, arg := range args {
		if phoneNumber, ok := arg.(string); ok {
			fmt.Printf("sending sms to %s,time:%v\n", phoneNumber, time.Now())
		}
	}
	return nil
}

func TestDelayQueue(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	host, _ := cache.Get("redis.host")
	password, _ := cache.Get("redis.password")
	dq := delayqueue.New(host.(string), "test-yasin",
		0, redis.WithPassWord(password.(string)))
	err := dq.Register(JobActionSMS{})
	if err != nil {
		log.Fatal(err)
	}
	dq.StartBackground()
	fmt.Printf("add job:%v\n", time.Now())
	err = dq.AddJob(delayqueue.DelayJob{
		ID:        (&JobActionSMS{}).ID(),
		DelayTime: 10,
		Args:      []interface{}{"181****9331"},
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Second)
}
