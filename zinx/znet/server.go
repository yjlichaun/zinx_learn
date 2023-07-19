package znet

import (
	"fmt"
	"net"
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
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
func (s *Server) Start() {
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
		fmt.Println("start Zinx server success,", s.Name, "success,listening ...")
		//循环监听
		for {
			//阻塞等待客户端链接，处理客户端链接业务（读写）
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp error :", err)
				continue
			}
			//已经与客户端建立链接，做一些业务，最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read tcp error :", err)
						continue
					}
					fmt.Printf("receive client buf %s , cnt = %d\n", buf, n)
					//回显功能
					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Println("write back buf error :", err)
						continue
					}
				}
			}()
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
