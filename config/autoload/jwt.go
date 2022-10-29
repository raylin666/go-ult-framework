package autoload

type JWT struct {
	App    string `yaml:"app"`
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}
