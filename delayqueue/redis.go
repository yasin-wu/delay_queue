package delayqueue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/yasin-wu/utils/redis"
)

type redisClient struct {
	keyPrefix  string
	batchLimit int
	client     *redis.Client
}

func (cli *redisClient) zadd(job DelayJob) error {
	key := cli.formatKey(job.ID)
	member, err := json.Marshal(job)
	if err != nil {
		return err
	}
	var z redis.Z
	z.Member = string(member)
	switch job.Type {
	case DelayTypeDuration:
		z.Score = job.DelayTime + time.Now().Unix()
	case DelayTypeDate:
		z.Score = job.DelayTime
	default:
		z.Score = job.DelayTime + time.Now().Unix()
	}
	return cli.client.ZAdd(key, z)
}

func (cli *redisClient) batchHandle(IDs []string) error {
	var wg = sync.WaitGroup{}
	wg.Add(len(IDs))
	for _, name := range IDs {
		key := cli.formatKey(name)
		go func(key string) {
			batch, lastScore, err := cli.getBatch(key)
			if err != nil {
				delayQueue.logger.ErrorF("get batch failed , error:%s", err.Error())
			} else {
				for _, item := range batch {
					var dj DelayJob
					if item.Member != "" {
						if err := json.Unmarshal([]byte(item.Member), &dj); err != nil {
							delayQueue.logger.ErrorF("json unmarshal failed , error:%s", err.Error())
							continue
						}
					}
					if executor, ok := delayQueue.jobExecutorFactory[dj.ID]; !ok {
						continue
					} else {
						if err := executor.action.Execute(dj.Args); err != nil {
							delayQueue.logger.ErrorF("job action execute failed , error:%s", err.Error())
						}
					}
				}
			}
			defer func() {
				if err != nil || len(batch) != 0 {
					if err := cli.clearBatch(key, lastScore); err != nil {
						delayQueue.logger.ErrorF("clear batch failed , error:%s", err.Error())
					}
				}
				wg.Done()
			}()
		}(key)
	}
	wg.Wait()
	return nil
}

func (cli *redisClient) getBatch(key string) ([]redis.Z, int64, error) {
	var redisZs []redis.Z
	var lastScore int64
	var err error
	redisZs, err = cli.client.ZRangeByScore(key, "0", fmt.Sprintf("%d", time.Now().Unix()),
		true, true, 0, cli.batchLimit)
	if err != nil || len(redisZs) == 0 {
		return redisZs, lastScore, err
	}
	lastScore = redisZs[len(redisZs)-1].Score
	redisZs, err = cli.client.ZRangeByScore(key, "0", fmt.Sprintf("%d", lastScore),
		true, true, 0, cli.batchLimit)
	return redisZs, lastScore, err
}

func (cli *redisClient) clearBatch(key string, lastScore int64) error {
	return cli.client.ZRemrangEByScore(key, "0", fmt.Sprintf("%d", lastScore))
}

func (cli *redisClient) formatKey(name string) string {
	return fmt.Sprintf("%s:%s", cli.keyPrefix, name)
}
