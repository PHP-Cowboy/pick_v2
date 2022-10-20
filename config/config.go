package config

type ServerConfig struct {
	Port       int            `mapstructure:"port" json:"port"`
	MysqlInfo  MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	JwtInfo    JWTConfig      `mapstructure:"jwt" json:"jwt"`
	RedisInfo  RedisConfig    `mapstructure:"redis" json:"redis"`
	GoodsApi   GoodsApiConfig `mapstructure:"goods_api" json:"goods_api"`
	U8Api      GoodsApiConfig `mapstructure:"u8_api" json:"u8_api"`
	YongSender string         `mapstructure:"yong_sender" json:"yong_sender"`
	RocketMQ   string         `mapstructure:"rocket_mq" json:"rocket_mq"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password"`
	Expire   int    `mapstructure:"expire" json:"expire"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type GoodsApiConfig struct {
	Url  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}
