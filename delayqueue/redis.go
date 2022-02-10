package delayqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisOptions redis.Options

type redisClient struct {
	keyPrefix  string
	batchLimit int64
	client     *redis.Client
	ctx        context.Context
}

var defaultRedisOptions = &RedisOptions{Addr: "localhost:6379", Password: "", DB: 0}

func (r *redisClient) zadd(job DelayJob) error {
	key := r.formatKey(job.ID)
	member, err := json.Marshal(job)
	if err != nil {
		return err
	}
	var z redis.Z
	z.Member = member
	switch job.Type {
	case DelayTypeDuration:
		z.Score = float64(job.DelayTime + time.Now().Unix())
	case DelayTypeDate:
		z.Score = float64(job.DelayTime)
	default:
		z.Score = float64(job.DelayTime + time.Now().Unix())
	}
	return r.client.ZAdd(r.ctx, key, &z).Err()
}

func (r *redisClient) batchHandle(IDs []string) error {
	var wg = sync.WaitGroup{}
	wg.Add(len(IDs))
	for _, name := range IDs {
		key := r.formatKey(name)
		go func(key string) {
			batch, lastScore, err := r.getBatch(key)
			if err != nil {
				delayQueue.logger.ErrorF("get batch failed , error:%s", err.Error())
			} else {
				for _, item := range batch {
					var dj DelayJob
					if item.Member != "" {
						if err := json.Unmarshal([]byte(item.Member.(string)), &dj); err != nil {
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
					if err := r.clearBatch(key, lastScore); err != nil {
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

func (r *redisClient) getBatch(key string) ([]redis.Z, int64, error) {
	var redisZs []redis.Z
	var lastScore int64
	var err error
	var opt redis.ZRangeBy
	opt.Min = "0"
	opt.Max = fmt.Sprintf("%d", time.Now().Unix())
	opt.Offset = 0
	opt.Count = r.batchLimit
	redisZs, err = r.client.ZRangeByScoreWithScores(r.ctx, key, &opt).Result()
	if err != nil || len(redisZs) == 0 {
		return redisZs, lastScore, err
	}
	lastScore = int64(redisZs[len(redisZs)-1].Score)
	opt.Max = fmt.Sprintf("%d", lastScore)
	redisZs, err = r.client.ZRangeByScoreWithScores(r.ctx, key, &opt).Result()
	return redisZs, lastScore, err
}

func (r *redisClient) clearBatch(key string, lastScore int64) error {
	return r.client.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("%d", lastScore)).Err()
}

func (r *redisClient) formatKey(name string) string {
	return fmt.Sprintf("%s:%s", r.keyPrefix, name)
}
