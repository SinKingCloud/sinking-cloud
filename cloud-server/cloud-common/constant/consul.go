package constant

// 云端配置
const (
	Postgres     = "postgres"         //postgres配置前缀
	PostgresHost = Postgres + ".host" //postgres地址
	PostgresPort = Postgres + ".port" //postgres端口
	PostgresUser = Postgres + ".user" //postgres用户
	PostgresPwd  = Postgres + ".pwd"  //postgres密码
	PostgresName = Postgres + ".name" //postgres数据库名称

	Redis     = "redis"         //redis配置前缀
	RedisHost = Redis + ".host" //redis地址
	RedisPort = Redis + ".port" //redis端口
	RedisPwd  = Redis + ".pwd"  //redis密码
	RedisDb   = Redis + ".db"   //redis编号
)
