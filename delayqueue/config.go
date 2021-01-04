package delayqueue

var (
	Setting *Config
)

const (
	// DefaultZRangeLimit
	DefaultZRangeLimit = 1000
	// DefaultBucketSize bucket数量
	DefaultBucketSize = 10
	// DefaultBucketName bucket名称
	DefaultBucketName = "dq_bucket_%d"
	// DefaultQueueName 队列名称
	DefaultQueueName = "dq_queue_%s"
	// DefaultQueueBlockTimeout 轮询队列超时时间 < Redis读取超时时间
	DefaultQueueBlockTimeout = 178
	// DefaultRedisHost Redis连接地址
	DefaultRedisHost = "127.0.0.1:6379"
	// DefaultRedisDb Redis数据库编号
	DefaultRedisDb = 1
	// DefaultRedisPassword Redis密码
	DefaultRedisPassword = ""
	// DefaultRedisMaxIdle Redis连接池闲置连接数
	DefaultRedisMaxIdle = 10
	// DefaultRedisMaxActive Redis连接池最大激活连接数, 0为不限制
	DefaultRedisMaxActive = 0
	// DefaultRedisConnectTimeout Redis连接超时时间,单位毫秒
	DefaultRedisConnectTimeout = 5000
	// DefaultRedisReadTimeout Redis读取超时时间, 单位毫秒
	DefaultRedisReadTimeout = 180000
	// DefaultRedisWriteTimeout Redis写入超时时间, 单位毫秒
	DefaultRedisWriteTimeout = 3000
	// DefaultRedisExpireTime Redis key 过期时间，单位秒
	DefaultRedisExpireTime = 86400
)

type Config struct {
	ZRangeLimit       int         // zrang max
	BucketSize        int         // bucket数量
	BucketName        string      // bucket在redis中的键名
	QueueName         string      // ready queue在redis中的键名
	QueueBlockTimeout int         // 调用blpop阻塞超时时间, 单位秒, 修改此项, redis.read_timeout必须做相应调整
	Redis             RedisConfig // redis配置
}

type RedisConfig struct {
	Host           string
	Db             int
	Password       string
	MaxIdle        int // 连接池最大空闲连接数
	MaxActive      int // 连接池最大激活连接数
	ConnectTimeout int // 连接超时,单位毫秒
	ReadTimeout    int // 读取超时,单位毫秒
	WriteTimeout   int // 写入超时,单位毫秒
	ExpireTime     int // key过期时间,单位秒
}

func (this *Config) initConfig(config *Config) {

	zrangeLimit := config.ZRangeLimit
	bucketSize := config.BucketSize
	bucketName := config.BucketName
	queueName := config.QueueName
	queueBlockTimeout := config.QueueBlockTimeout

	redisHost := config.Redis.Host
	redisDb := config.Redis.Db
	redisPassword := config.Redis.Password
	redisMaxIdle := config.Redis.MaxIdle
	redisMaxActive := config.Redis.MaxActive
	redisConnectTimeout := config.Redis.ConnectTimeout
	redisReadTimeout := config.Redis.ReadTimeout
	redisWriteTimeout := config.Redis.WriteTimeout
	redisExpireTime := config.Redis.ExpireTime

	if queueBlockTimeout > redisReadTimeout {
		queueBlockTimeout = redisReadTimeout/100 - 2
	}
	if zrangeLimit == 0 {
		zrangeLimit = DefaultZRangeLimit
	}
	if bucketSize == 0 {
		bucketSize = DefaultBucketSize
	}
	if bucketName == "" {
		bucketName = DefaultBucketName
	}
	if queueName == "" {
		queueName = DefaultQueueName
	}
	if queueBlockTimeout == 0 {
		queueBlockTimeout = DefaultQueueBlockTimeout
	}
	if redisHost == "" {
		redisHost = DefaultRedisHost
	}
	if redisMaxIdle == 0 {
		redisMaxIdle = DefaultRedisMaxIdle
	}
	if redisConnectTimeout == 0 {
		redisConnectTimeout = DefaultRedisConnectTimeout
	}
	if redisReadTimeout == 0 {
		redisReadTimeout = DefaultRedisReadTimeout
	}
	if redisWriteTimeout == 0 {
		redisWriteTimeout = DefaultRedisWriteTimeout
	}
	if redisExpireTime == 0 {
		redisExpireTime = DefaultRedisExpireTime
	}
	this.ZRangeLimit = zrangeLimit
	this.BucketSize = bucketSize
	this.BucketName = bucketName
	this.QueueName = queueName
	this.QueueBlockTimeout = queueBlockTimeout

	this.Redis.Host = redisHost
	this.Redis.Db = redisDb
	this.Redis.Password = redisPassword
	this.Redis.MaxIdle = redisMaxIdle
	this.Redis.MaxActive = redisMaxActive
	this.Redis.ConnectTimeout = redisConnectTimeout
	this.Redis.ReadTimeout = redisReadTimeout
	this.Redis.WriteTimeout = redisWriteTimeout
	this.Redis.ExpireTime = redisExpireTime
}

// 初始化默认配置
func (this *Config) initDefaultConfig() {
	this.ZRangeLimit = DefaultZRangeLimit
	this.BucketSize = DefaultBucketSize
	this.BucketName = DefaultBucketName
	this.QueueName = DefaultQueueName
	this.QueueBlockTimeout = DefaultQueueBlockTimeout

	this.Redis.Host = DefaultRedisHost
	this.Redis.Db = DefaultRedisDb
	this.Redis.Password = DefaultRedisPassword
	this.Redis.MaxIdle = DefaultRedisMaxIdle
	this.Redis.MaxActive = DefaultRedisMaxActive
	this.Redis.ConnectTimeout = DefaultRedisConnectTimeout
	this.Redis.ReadTimeout = DefaultRedisReadTimeout
	this.Redis.WriteTimeout = DefaultRedisWriteTimeout
	this.Redis.ExpireTime = DefaultRedisExpireTime
}
