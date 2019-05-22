package main

import (
	"gopcp.v2/chapter7/lottery/bootstrap"
	"gopcp.v2/chapter7/lottery/web/middleware/identity"
	"gopcp.v2/chapter7/lottery/web/routes"
	"fmt"
	"gopcp.v2/chapter7/lottery/conf"
	)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New("Go抽奖系统", "zchd")
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)

	return app
}

func main() {
	// 服务器集群的时候才需要区分这项设置
	// 比如：根据服务器的IP、名称、端口号等，或者运行的参数
	if port == 8080 {
		conf.RunningCrontabService = true
	}

	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
