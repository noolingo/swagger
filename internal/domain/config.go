package domain

type Config struct {
	App    app    `yaml:"app" env-prefix:"XPRTS_" json:"app"`
	Listen listen `yaml:"listen" env-prefix:"XPRTS_LISTEN_" json:"listen"`
	Log    Logger `yaml:"log" env-prefix:"XPRTS_" json:"log"`
}

type app struct {
	FilesPath string `yaml:"filesPath"`
	URL       string `yaml:"url" env:"BACKEND_URL"`
}

type listen struct {
	Ports ports  `yaml:"ports" env-prefix:"PORTS_"`
	Host  string `yaml:"host" env:"HOSTNAME" env-default:"0.0.0.0"`
}

type ports struct {
	HTTP string `yaml:"http" env:"HTTP"`
}

type Logger struct {
	Level map[string]string `yaml:"level" env:"LOG_LEVEL"`
}
