package admin

type Controller struct {
	Balance balanceController
	Order   orderController
	Log     logController
	Channel channelController
	Config  configController
}
