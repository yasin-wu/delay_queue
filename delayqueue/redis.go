package delayqueue

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

// 初始化连接池
func initRedisPool() {
	RedisPool = &redis.Pool{
		MaxIdle:      Setting.Redis.MaxIdle,
		MaxActive:    Setting.Redis.MaxActive,
		IdleTimeout:  300 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisTestOnBorrow,
		Wait:         true,
	}
}

func redisDial() (redis.Conn, error) {
	conn, err := redis.Dial(
		"tcp",
		Setting.Redis.Host,
		redis.DialConnectTimeout(time.Duration(Setting.Redis.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(Setting.Redis.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(Setting.Redis.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Printf("连接redis失败#%s", err.Error())
		return nil, err
	}

	if Setting.Redis.PassWord != "" {
		if _, err := conn.Do("AUTH", Setting.Redis.PassWord); err != nil {
			conn.Close()
			log.Printf("redis认证失败#%s", err.Error())
			return nil, err
		}
	}

	_, err = conn.Do("SELECT", Setting.Redis.Db)
	if err != nil {
		conn.Close()
		log.Printf("redis选择数据库失败#%s", err.Error())
		return nil, err
	}

	return conn, nil
}

func redisTestOnBorrow(conn redis.Conn, t time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		log.Printf("从redis连接池取出的连接无效#%s", err.Error())
	}
	return err
}

func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := RedisPool.Get()
	defer redis.Close()
	return redis.Do(command, args...)
}
