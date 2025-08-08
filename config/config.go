package config

type Config struct {
	HttpServer HttpServerConfig
	Redis      RedisConfig
}

type HttpServerConfig struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

type RedisConfig struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}
