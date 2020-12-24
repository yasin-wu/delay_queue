package delayqueue

import (
	"fmt"
	"time"
)

var (
	// 每个定时器对应一个bucket
	timers []*time.Ticker
	// bucket名称chan
	bucketNameChan <-chan string
)

// Init 初始化配置
func (this *DelayCmd) InitConfig() {
	Setting = &Config{}
	if this.Config == nil {
		Setting.initDefaultConfig()
		return
	}
	Setting.initConfig(this.Config)
}

// Init 初始化延时队列
func (this *DelayCmd) Init() {
	initRedisPool()
	initTimers()
	bucketNameChan = generateBucketName()
}

func initTimers() {
	timers = make([]*time.Ticker, Setting.BucketSize)
	var bucketName string
	for i := 0; i < Setting.BucketSize; i++ {
		timers[i] = time.NewTicker(1 * time.Second)
		bucketName = fmt.Sprintf(Setting.BucketName, i+1)
		go waitTicker(timers[i], bucketName)
	}
}

func generateBucketName() <-chan string {
	c := make(chan string)
	go func() {
		i := 1
		for {
			c <- fmt.Sprintf(Setting.BucketName, i)
			if i >= Setting.BucketSize {
				i = 1
			} else {
				i++
			}
		}
	}()
	return c
}
