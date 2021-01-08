package delayqueue

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yasin-wu/delay-queue/cronjob"
	"github.com/yasin-wu/delay-queue/logger"
)

func initDelayQueue(conf *Config) {
	delayQueue = new(DelayQueue)
	checkConf(conf)
	sche := cronjob.New()
	sche.Register([]int{}, 1, DelayQueueCronJob{})
	delayQueue.logger = logger.DefaultLogger
	delayQueue.scheduler = sche
	delayQueue.jobExecutorFactory = make(map[string]*jobExecutor)
	delayQueue.redisCli = initRedis(conf)
	delayQueue.redisConf = conf.Redis
	delayQueue.logger.InfoF("Initialization of DelayQueue completed......")
}

func initRedis(conf *Config) *redisClient {
	redisPool = &redis.Pool{
		MaxIdle:      conf.Redis.MaxIdle,
		MaxActive:    conf.Redis.MaxActive,
		IdleTimeout:  defaultIdleTimeout,
		Dial:         redisDial,
		TestOnBorrow: redisTestOnBorrow,
		Wait:         true,
	}
	redisCli = &redisClient{
		keyPrefix:  conf.KeyPrefix,
		batchLimit: conf.BatchLimit,
	}
	return redisCli
}

func checkConf(conf *Config) {
	if conf.KeyPrefix == "" {
		conf.KeyPrefix = defaultKeyPrefix
	}
	if conf.BatchLimit == 0 {
		conf.BatchLimit = defaultBatchLimit
	}

	if conf.Redis.Host == "" {
		conf.Redis.Host = defaultRedisHost
	}
	if conf.Redis.PassWord == "" {
		conf.Redis.PassWord = defaultRedisPassWord
	}
	if conf.Redis.DB == 0 {
		conf.Redis.DB = defaultRedisDB
	}
	if conf.Redis.MaxIdle == 0 {
		conf.Redis.MaxIdle = defaultRedisMaxIdle
	}
	if conf.Redis.MaxActive == 0 {
		conf.Redis.MaxActive = defaultRedisMaxActive
	}
	if conf.Redis.ConnectTimeout == 0 {
		conf.Redis.ConnectTimeout = defaultRedisConnectTimeout
	}
	if conf.Redis.ReadTimeout == 0 {
		conf.Redis.ReadTimeout = defaultRedisReadTimeout
	}
	if conf.Redis.WriteTimeout == 0 {
		conf.Redis.WriteTimeout = defaultRedisWriteTimeout
	}
}

func redisDial() (redis.Conn, error) {
	conn, err := redis.Dial(
		"tcp",
		delayQueue.redisConf.Host,
		redis.DialConnectTimeout(time.Duration(delayQueue.redisConf.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(delayQueue.redisConf.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(delayQueue.redisConf.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		delayQueue.logger.ErrorF("connect redis failed , error:%s", err.Error())
		return nil, err
	}

	if delayQueue.redisConf.PassWord != "" {
		if _, err := conn.Do("AUTH", delayQueue.redisConf.PassWord); err != nil {
			conn.Close()
			delayQueue.logger.ErrorF("auth redis failed , error:%s", err.Error())
			return nil, err
		}
	}

	_, err = conn.Do("SELECT", delayQueue.redisConf.DB)
	if err != nil {
		conn.Close()
		delayQueue.logger.ErrorF("select redis db failed , error:%s", err.Error())
		return nil, err
	}

	return conn, nil
}

func redisTestOnBorrow(conn redis.Conn, t time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		delayQueue.logger.ErrorF("get redis conn from redis pool failed , error:%s", err.Error())
	}
	return err
}
