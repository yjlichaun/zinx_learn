package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	//服务器名称
	Name string
	//服务器ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务监听的端口
	Port int
	//当前server添加router server注册的链接对应的处理业务
	Router ziface.IRouter
}

//CallBackToClient 定义当前客户端链接所绑定的handle api (目前是写死的，后期应该优化有用户自定的handle方法)
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	//回显业务
//	fmt.Println("[Conn Handle] CallbackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf error", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
}
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listener at Ip : %s , Port : %d is starting", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version : %s , MaxConn : %d , MaxPacketSize : %d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize,
	)
	fmt.Println("[start] Server Listener at Ip :", s.IP+", port :", s.Port, "is starting")
	go func() {
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
			return
		}
		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp error :", err)
			return
		}
		var connId uint32
		connId = 0
		fmt.Println("start Zinx server success,", s.Name, "success,listening ...")
		//循环监听
		for {
			//阻塞等待客户端链接，处理客户端链接业务（读写）
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp error :", err)
				continue
			}
			//绑定链接和业务，得到连接模块
			dealConn := NewConnection(conn, connId, s.Router)
			connId++
			//启动当前链接处理业务
			go dealConn.Start()
		}

	}()

}
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息，进行停止或者回收
}
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务
	//阻塞主函数
	select {}
}

//AddRouter 添加一个路由
func (s *Server) AddRouter(route ziface.IRouter) {
	s.Router = route
	fmt.Println("add router success !!!")
}
