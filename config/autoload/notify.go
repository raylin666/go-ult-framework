package autoload

type Notify struct {
	Recover struct{
		Email struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
			To   string `yaml:"to"`
		} `yaml:"email"`
	} `yaml:"recover"`
}
