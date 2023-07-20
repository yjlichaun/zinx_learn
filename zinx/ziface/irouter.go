package ziface

// IRouter 路由抽象接口
//路由里的数据都是IRequest
type IRouter interface {
	//PreHandle 处理Conn业务之前的hook
	PreHandle(req IRequest)
	//Handle 处理Conn业务的hook
	Handle(req IRequest)
	//PostHandle 处理Conn业务之后的hook
	PostHandle(req IRequest)
}
