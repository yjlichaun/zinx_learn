package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

//连接模块
type Connection struct {
	//socket TcpSocket
	Conn *net.TCPConn
	//链接id
	ConnId uint32
	//链接status (是否关闭)
	ConnStatus bool
	//等待连接被动退出的channel
	ExitChan chan bool
	//该链接处理的方法router
	Router ziface.IRouter
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:       conn,
		ConnId:     connId,
		ConnStatus: false,
		Router:     router,
		ExitChan:   make(chan bool, 1),
	}
}

// StartReader 链接的读业务方法
func (conn *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("connId = ", conn.ConnId, "Reader is exit, remote addr is ,", conn.GetRemoteAddr().String())
	defer conn.Stop()
	for {
		buf := make([]byte, 1024)
		_, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("conn Read error: ", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: conn,
			data: buf,
		}

		//执行注册的路由方法
		go func(req ziface.IRequest) {
			conn.Router.PreHandle(req)
			conn.Router.Handle(req)
			conn.Router.PostHandle(req)
		}(&req)
		////调用当前链接绑定的HandleApi
		//if err := conn.HandleApi(conn.Conn, buf, n); err != nil {
		//	fmt.Println("connId ", conn.ConnId, "HandleApi error: ", err)
		//	break
		//}
	}
}

//实现方法 ----------------------------------------------------------------
func (conn *Connection) Start() {
	fmt.Println("conn started ... connId:", conn.ConnId)
	//启动当前链接的读数据业务
	conn.StartReader()
	//TODO: 启动当前链接的写数据业务

}
func (conn *Connection) Stop() {
	fmt.Println("conn stop() .. ConnId: ", conn.ConnId)
	if conn.ConnStatus == true {
		return
	}
	conn.ConnStatus = true
	//关闭链接
	conn.Conn.Close()
	//关闭管道
	close(conn.ExitChan)
}
func (conn *Connection) GetTcpConnection() *net.TCPConn {
	return conn.Conn
}
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnId
}
func (conn *Connection) GetRemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}
func (conn *Connection) Send(data []byte) error {
	return nil
}
