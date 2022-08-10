[![OSCS Status](https://www.oscs1024.com/platform/badge/yasin-wu/delay_queue.svg?size=small)](https://www.murphysec.com/dr/kFJ0vHLhJQTz8wiubq)
## 介绍

delay queue 是基于Redis Zset实现的Golang版延时队列。以时间戳作为Score, 主动轮询小于当前时间的元素。新增延迟类型支持：支持延迟多少秒和延迟到具体时间。

## 安装

```
go get -u github.com/yasin-wu/delay_queue
```

推荐使用go.mod

```
require github.com/yasin-wu/delay_queue/v2 v2.1.0
```

## 使用

```go
var redisOptions = &delayqueue.RedisOptions{Addr: "47.108.155.25:6379", Password: "yasinwu"}

type JobActionSMS struct{}

var _ pkg.JobBaseAction = (*JobActionSMS)(nil)

func (JobActionSMS) ID() string {
    return "JobActionSMS"
}

func (JobActionSMS) Execute(args []any) error {
    for _, arg := range args {
        if phoneNumber, ok := arg.(string); ok {
            fmt.Printf("sending sms to %s,time:%v\n", phoneNumber, time.Now())
        }
    }
    return nil
}

func main() {
    dq := dqueue.New("test-yasin", 0, redisOptions)
    err := dq.Register(JobActionSMS{})
    if err != nil {
        log.Fatal(err)
    }
    dq.StartBackground()
    fmt.Printf("add job:%v\n", time.Now())
    err = dq.AddJob(pkg.DelayJob{
        ID:        (&JobActionSMS{}).ID(),
        DelayTime: 10,
        Args:      []any{"181****9331"},
    })
    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(20 * time.Second)
}

```
