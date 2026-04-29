package constant

const (
	LockBalanceOperate     = "BalanceOperate_" //余额操作并发锁
	LockTimeBalanceOperate = 1 * 60 * 1000     //余额操作并发锁最大时间
	LockOrderCallback      = "OrderCallback_"  //订单回调并发锁
	LockTimeOrderCallback  = 1 * 60 * 1000     //订单回调并发锁最大时间

	LockOrderCleanExpiredUnPay     = "OrderCleanExpiredUnPay" //清理过期未支付订单锁
	LockTimeOrderCleanExpiredUnPay = 10 * 60 * 1000           //清理过期未支付订单锁最大时间
)
