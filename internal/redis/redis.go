package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yasin-wu/delay_queue/v2/internal/logger"

	"github.com/yasin-wu/delay_queue/v2/pkg"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	keyPrefix  string
	batchLimit int64
	client     *redis.Client
	ctx        context.Context
	logger     logger.Logger
}

func New(keyPrefix string, batchLimit int64, redisOptions *pkg.Options) *Client {
	return &Client{
		keyPrefix:  keyPrefix,
		batchLimit: batchLimit,
		client:     redis.NewClient((*redis.Options)(redisOptions)),
		ctx:        context.Background(),
		logger:     logger.DefaultLogger,
	}
}

func (c *Client) Zadd(job pkg.DelayJob) error {
	key := c.FormatKey(job.ID)
	delayTime := job.DelayTime
	job.DelayTime = -1
	member, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	var z redis.Z
	z.Member = member
	z.Score = float64(delayTime + time.Now().Unix())
	if job.Type == pkg.DelayTypeDate {
		z.Score = float64(delayTime)
	}
	return c.client.ZAdd(c.ctx, key, &z).Err()
}

func (c *Client) Zremove(job pkg.DelayJob) error {
	key := c.FormatKey(job.ID)
	job.DelayTime = -1
	member, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	return c.client.ZRem(c.ctx, key, member).Err()
}

func (c *Client) GetBatch(key string) ([]redis.Z, int64, error) {
	var redisZs []redis.Z
	var lastScore int64
	var err error
	var opt redis.ZRangeBy
	opt.Min = "0"
	opt.Max = fmt.Sprintf("%d", time.Now().Unix())
	opt.Offset = 0
	opt.Count = c.batchLimit
	redisZs, err = c.client.ZRangeByScoreWithScores(c.ctx, key, &opt).Result()
	if len(redisZs) > 0 {
		lastScore = int64(redisZs[len(redisZs)-1].Score)
	}
	//if err != nil || len(redisZs) == 0 {
	//	return redisZs, lastScore, err
	//}
	//lastScore = int64(redisZs[len(redisZs)-1].Score)
	//opt.Max = fmt.Sprintf("%d", lastScore)
	//redisZs, err = c.client.ZRangeByScoreWithScores(c.ctx, key, &opt).Result()
	return redisZs, lastScore, err
}

func (c *Client) ClearBatch(key string, lastScore int64) {
	if err := c.client.ZRemRangeByScore(c.ctx, key, "0", fmt.Sprintf("%d", lastScore)).Err(); err != nil {
		c.logger.Errorf("clear batch failed , error:%s", err.Error())
	}
}

func (c *Client) FormatKey(name string) string {
	return fmt.Sprintf("%s:%s", c.keyPrefix, name)
}

func (c *Client) SetLogger(logger logger.Logger) {
	c.logger = logger
}
