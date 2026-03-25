package constant

// 云端配置
const (
	Postgres     = "postgres"         //postgres 配置前缀
	PostgresHost = Postgres + ".host" //postgres 地址
	PostgresPort = Postgres + ".port" //postgres 端口
	PostgresUser = Postgres + ".user" //postgres 用户
	PostgresPwd  = Postgres + ".pwd"  //postgres 密码
	PostgresName = Postgres + ".name" //postgres 数据库名称

	Redis     = "redis"         //redis 配置前缀
	RedisHost = Redis + ".host" //redis 地址
	RedisPort = Redis + ".port" //redis 端口
	RedisPwd  = Redis + ".pwd"  //redis 密码
	RedisDb   = Redis + ".db"   //redis 编号

	Nats     = "nats"
	NatsHost = Nats + ".host" //nats 地址
	NatsPort = Nats + ".port" //nats 端口
	NatsUser = Nats + ".user" //nats 用户
	NatsPwd  = Nats + ".pwd"  //nats 密码
)
