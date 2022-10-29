package autoload

type DB struct {
	Dsn         string `yaml:"dsn"`
	Driver      string `yaml:"driver"`
	DbName      string `yaml:"db_name"`
	Host        string `yaml:"host"`
	UserName    string `yaml:"user_name"`
	Password    string `yaml:"password"`
	Charset     string `yaml:"charset"`
	Port        int    `yaml:"port"`
	Prefix      string `yaml:"prefix"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxLifeTime int64  `yaml:"max_life_time"`
	ParseTime   string `yaml:"parse_time"`
	Loc         string `yaml:"loc"`
}
