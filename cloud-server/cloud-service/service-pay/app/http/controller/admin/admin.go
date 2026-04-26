package admin

type Controller struct {
	Balance balanceController
	Order   orderController
	Channel channelController
	Config  configController
}
