package autoload

type Server struct {
	Http Http
}

type Http struct {
	Network string `yaml:"network"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Cors    struct {
		Domains string `yaml:"domains"`
	} `yaml:"cors"`
}
