package restapi

type ServerConfig struct {
	Mongo struct {
		ConnectionString string `yaml:"connection-string"`
		DatabaseName     string `yaml:"db-name"`
		User             string `yaml:"user"`
		Password         string `yaml:"password"`
	} `yaml:"mongo"`
}
