package config

type Config struct {
	LogLevel   string           `mapstructure:"LOG_LEVEL"`
	GRPCServer GRPCServerConfig `mapstructure:"GRPC_SERVER"`
	Client     ClientConfig     `mapstructure:"CLIENT"`
}

type GRPCServerConfig struct {
	Port int `mapstructure:"PORT"`
}

type ClientConfig struct {
	ServerList []string `mapstructure:"SERVER_LIST"`
	Timeout    string   `mapstructure:"TIMEOUT"`
	ChunkSize  int      `mapstructure:"CHUNK_SIZE"`
}
