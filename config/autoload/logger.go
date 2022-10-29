package autoload

type Logger struct {
	MaxSize    int  `yaml:"max_size"`
	MaxBackups int  `yaml:"max_backups"`
	MaxAge     int  `yaml:"max_age"`
	LocalTime  bool `yaml:"local_time"`
	Compress   bool `yaml:"compress"`
}
