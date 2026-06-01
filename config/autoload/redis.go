package autoload

// Redis Redis 连接配置结构体。
type Redis struct {
	Network         string `yaml:"network"`           // 网络协议 (tcp)
	Addr            string `yaml:"addr"`              // Redis 地址
	Port            int    `yaml:"port"`              // Redis 端口
	Username        string `yaml:"username"`          // 用户名
	Password        string `yaml:"password"`          // 密码
	DB              int    `yaml:"db"`                // 数据库索引
	MaxRetries      int    `yaml:"max_retries"`       // 连接最大重试次数
	RetryDelay      int    `yaml:"retry_delay"`       // 重试间隔（秒）
	MinRetryBackoff int64  `yaml:"min_retry_backoff"` // 最小重试退避时间
	MaxRetryBackoff int64  `yaml:"max_retry_backoff"` // 最大重试退避时间
	DialTimeout     int64  `yaml:"dial_timeout"`      // 连接超时时间
	ReadTimeout     int64  `yaml:"read_timeout"`      // 读超时时间
	WriteTimeout    int64  `yaml:"write_timeout"`     // 写超时时间
	PoolFIFO        bool   `yaml:"pool_fifo"`         // 连接池 FIFO 模式
	PoolSize        int    `yaml:"pool_size"`         // 连接池大小
	MinIdleConns    int    `yaml:"min_idle_conns"`    // 最小空闲连接数
	MaxConnAge      int64  `yaml:"max_conn_age"`      // 连接最大存活时间
	PoolTimeout     int64  `yaml:"pool_timeout"`      // 连接池超时时间
	IdleTimeout     int64  `yaml:"idle_timeout"`      // 空闲连接超时时间
}
