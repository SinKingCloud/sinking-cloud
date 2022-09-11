package constant

const (
	LockWebUserReg     = "WebUserReg_" //注册用户并发锁
	LockTimeWebUserReg = 1 * 60 * 1000 //注册用户并发锁最大时间

	LockChangeOrderStatus     = "Order_"      //订单状态更改并发锁
	LockTimeChangeOrderStatus = 1 * 60 * 1000 //订单状态更改并发锁最大时间
)
