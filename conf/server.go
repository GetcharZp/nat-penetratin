package conf

type Server struct {
	Web struct {
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"web"`
}
