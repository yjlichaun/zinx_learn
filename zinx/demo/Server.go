package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//ping rest 自定义路由
type PingRouter struct {
	znet.BaseRouter
}
type HelloZinxRouter struct {
	znet.BaseRouter
}

//func (p *PingRouter) PreHandle(req ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before ping... \n"))
//	if err != nil {
//		fmt.Println("call back before ping error", err)
//	}
//}

func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//read client data ,write ping...
	fmt.Println("recv from client :msgId = ", req.GetMsgId(), ", data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(200, []byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping error", err)
	}
}
func (p *HelloZinxRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	//read client data ,write ping...
	fmt.Println("recv from client :msgId = ", req.GetMsgId(), ", data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(201, []byte("hello...\n"))
	if err != nil {
		fmt.Println("call back ping error", err)
	}

}

//func (p *PingRouter) PostHandle(req ziface.IRequest) {
//	fmt.Println("Call Router PostHandle...")
//	_, err := req.GetConnection().GetTcpConnection().Write([]byte("After ping... \n"))
//	if err != nil {
//		fmt.Println("call back after ping error", err)
//	}
//}

func main() {
	z := znet.NewServer("[demoServer]")
	//给当前框架添加router
	z.AddRouter(0, &PingRouter{})
	z.AddRouter(1, &HelloZinxRouter{})
	z.Serve()
}
