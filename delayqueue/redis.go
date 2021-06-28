package delayqueue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisZ struct {
	Score  int64
	Member string
}

type redisClient struct {
	keyPrefix  string
	batchLimit int64
}

func (cli *redisClient) ZAdd(job DelayJob) error {
	key := cli.formatKey(job.ID)
	member, err := json.Marshal(job)
	if err != nil {
		return err
	}
	switch job.Type {
	case DelayTypeDuration:
		_, err = cli.execRedisCommand("ZADD", key, job.DelayTime+time.Now().Unix(), member)
	case DelayTypeDate:
		_, err = cli.execRedisCommand("ZADD", key, job.DelayTime, member)
	default:
		_, err = cli.execRedisCommand("ZADD", key, job.DelayTime+time.Now().Unix(), member)
	}
	return err
}

func (cli *redisClient) BatchHandle(IDs []string) error {
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

func (cli *redisClient) getBatch(key string) ([]RedisZ, int64, error) {
	var redisZs []RedisZ
	var lastScore int64
	var err error
	batchVal, err := redis.Values(
		cli.execRedisCommand("ZRANGEBYSCORE", key,
			0, time.Now().Unix(),
			"WITHSCORES",
			"limit", 0, cli.batchLimit))
	if err != nil || len(batchVal) == 0 {
		return redisZs, lastScore, err
	}
	redisZs = cli.handleBatchVal(batchVal)
	lastScore = redisZs[len(redisZs)-1].Score
	batchVal, err = redis.Values(cli.execRedisCommand("ZRANGEBYSCORE", key,
		0, lastScore,
		"WITHSCORES",
		"limit", 0, cli.batchLimit))
	redisZs = cli.handleBatchVal(batchVal)
	return redisZs, lastScore, err
}

func (cli *redisClient) handleBatchVal(batchVal []interface{}) []RedisZ {
	var err error
	redisZs := make([]RedisZ, len(batchVal)/2)
	for i := 0; i < len(redisZs); i++ {
		var redisZ RedisZ
		redisZ.Member = readString(batchVal[i*2])
		redisZ.Score, err = strconv.ParseInt(readString(batchVal[i*2+1]), 10, 64)
		if err != nil {
			delayQueue.logger.ErrorF("string to int64 failed , error:%s", err.Error())
		}
		redisZs[i] = redisZ
	}
	return redisZs
}

func readString(value interface{}) string {
	var buffer []byte
	typeString := reflect.TypeOf(value).String()
	switch typeString {
	case "[]uint8":
		for _, v := range value.([]uint8) {
			buffer = append(buffer, v)
		}
	}
	return string(buffer)
}

func (cli *redisClient) clearBatch(key string, lastScore int64) error {
	_, err := cli.execRedisCommand("ZREMRANGEBYSCORE", key, 0, lastScore)
	return err
}

func (cli *redisClient) formatKey(name string) string {
	return fmt.Sprintf("%s:%s", cli.keyPrefix, name)
}

func (cli *redisClient) execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := redisPool.Get()
	defer redis.Close()
	return redis.Do(command, args...)
}
