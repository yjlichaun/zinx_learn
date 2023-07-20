package ziface

type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	//AddRouter 路由功能
	AddRouter(msgId uint32, router IRouter)
}
