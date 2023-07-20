package znet

import "zinx/ziface"

//BaseRouter 实现router时,先嵌入这个BaseRoute基类，然后根据需求对这个基类进行方法重写
type BaseRouter struct {
}

//PreHandle 处理Conn业务之前的hook
func (r *BaseRouter) PreHandle(req ziface.IRequest) {}

//Handle 处理Conn业务的hook
func (r *BaseRouter) Handle(req ziface.IRequest) {}

//PostHandle 处理Conn业务之后的hook
func (r *BaseRouter) PostHandle(req ziface.IRequest) {}
