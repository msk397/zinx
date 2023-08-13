package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()

	// GetConnMgr 获取ConnMgr
	GetConnMgr() IConnmanager
	// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)

	// SetOnConnStart 注册OnConnStart钩子函数的方法
	SetOnConnStart(func(connection IConnection))

	// SetOnConnStop 注册OnConnStop钩子函数的方法
	SetOnConnStop(func(connection IConnection))

	// CallOnConnStart 调用OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)

	// CallOnConnStop 调用OnConnStop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
