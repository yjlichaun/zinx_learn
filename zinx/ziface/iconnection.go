package ziface

import "net"

type IConnection interface {
	//启动链接 让当前连接开始工作
	Start()
	//停止链接 让当前连接停止工作
	Stop()
	//获取当前链接的绑定socket connection
	GetTcpConnection() *net.TCPConn
	//获取当前链接模块的链接id
	GetConnID() uint32
	//获取远程客户端的tcp状态 ip port
	GetRemoteAddr() net.Addr
	//发送数据，将数据发送给远程客户端
	Send(data []byte) error
}

//HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
