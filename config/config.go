package config

type ServerConfig struct {
	Port      int         `mapstructure:"port" json:"port"`
	MysqlInfo MysqlConfig `mapstructure:"mysql" json:"mysql"`
	JwtInfo   JWTConfig   `mapstructure:"jwt" json:"jwt"`
	RedisInfo RedisConfig `mapstructure:"redis" json:"redis"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}
