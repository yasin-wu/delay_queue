package delayqueue

type Config struct {
	KeyPrefix  string
	BatchLimit int64
	Redis      *RedisConf
}

type RedisConf struct {
	Host           string
	PassWord       string
	DB             int
	MaxIdle        int // 连接池最大空闲连接数
	MaxActive      int // 连接池最大激活连接数
	ConnectTimeout int // 连接超时,单位毫秒
	ReadTimeout    int // 读取超时,单位毫秒
	WriteTimeout   int // 写入超时,单位毫秒
}
