package config

//type Config struct {
//	HttpServer HttpServerConfig `koanf:"http"`
//	Redis      RedisConfig      `koanf:"redis"`
//}

type Config struct {
	Http struct {
		Port string `koanf:"port"`
	} `koanf:"http"`
	Redis struct {
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		Password string `koanf:"password"`
	} `koanf:"redis"`
}

//type HttpServerConfig struct {
//	Port string `koanf:"port"`
//}
//
//type RedisConfig struct {
//	Addr     string `koanf:"addr"`
//	Port     string `koanf:"port"`
//	Password string `koanf:"password"`
//}
