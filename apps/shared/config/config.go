package config

type MongoDbConfig struct {
	Url string
}

type Config struct {
	MongoDb MongoDbConfig
}

func NewConfigFromEnvVars() *Config {
	return &Config{
		MongoDb: MongoDbConfig{
			Url: "JA",
		},
	}
}
