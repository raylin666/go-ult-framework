package autoload

type Redis struct {
	Network            string `yaml:"network"`
	Addr               string `yaml:"addr"`
	Port               int    `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	DB                 int    `yaml:"db"`
	MaxRetries         int    `yaml:"max_retries"`
	MinRetryBackoff    int64  `yaml:"min_retry_backoff"`
	MaxRetryBackoff    int64  `yaml:"max_retry_backoff"`
	DialTimeout        int64  `yaml:"dial_timeout"`
	ReadTimeout        int64  `yaml:"read_timeout"`
	WriteTimeout       int64  `yaml:"write_timeout"`
	PoolFIFO           bool   `yaml:"pool_fifo"`
	PoolSize           int    `yaml:"pool_size"`
	MinIdleConns       int    `yaml:"min_idle_conns"`
	MaxConnAge         int64  `yaml:"max_conn_age"`
	PoolTimeout        int64  `yaml:"pool_timeout"`
	IdleTimeout        int64  `yaml:"idle_timeout"`
	IdleCheckFrequency int64  `yaml:"idle_check_frequency"`
}
