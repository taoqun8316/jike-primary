package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(127.0.0.1:3306)/jike?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "127.0.0.1:6379",
	},
}
