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
	//GetConnMgr :get conn manager
	GetConnMgr() IConnManager
	//SetOnConnStart :func to register OnConnStart hook function
	SetOnConnStart(hookFunc func(conn IConnection))
	//SetOnConnStop :func to register OnConnStop hook function
	SetOnConnStop(hookFunc func(conn IConnection))
	//CallOnConnStart func to execute OnConnStart hook function
	CallOnConnStart(conn IConnection)
	//CallOnConnStop func to execute OnConnStop hook function
	CallOnConnStop(conn IConnection)
}
