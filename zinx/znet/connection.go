package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
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
	//The current server message management module is used to bind msgID and the corresponding processing business API relationship
	MsgHandler ziface.IMsgHandler
	//Unbuffered channel for message communication between reads and writes, goroutines
	MsgChan chan []byte
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnId:     connId,
		ConnStatus: false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
	}
}

// StartReader 链接的读业务方法
func (conn *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connId = ", conn.ConnId, "[Reader is exit]  , remote addr is ,", conn.GetRemoteAddr().String())
	defer conn.Stop()
	for {
		//buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		//_, err := conn.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("conn Read error: ", err)
		//	continue
		//}
		//create an object to pack package and unpack package
		dp := NewDataPack()
		//get client msg Head binary steam 8 bytes
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn.Conn, headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}

		//unpack get msgId and msgDataLen ->msg
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack msg error: ", err)
			break
		}
		//with dataLen read msgData second -> msgData
		var msgData []byte
		if msg.GetMsgLen() > 0 {
			msgData = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn.Conn, msgData); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(msgData)
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: conn,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			conn.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go conn.MsgHandler.DoMsgHandler(&req)
		}
		//执行注册的路由方法
		//go func(req ziface.IRequest) {
		//	conn.Router.PreHandle(req)
		//	conn.Router.Handle(req)
		//	conn.Router.PostHandle(req)
		//}(&req)

		//Call the HandleHandle API of the current link binding
		//Find the corresponding processing business execution according to the bound msgId
		//go conn.MsgHandler.DoMsgHandler(&req)

		//if err := conn.HandleApi(conn.Conn, buf, n); err != nil {
		//	fmt.Println("connId ", conn.ConnId, "HandleApi error: ", err)
		//	break
		//}
	}
}

//StartWriter The linked write business method, which is used to send data to the client
func (conn *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("connId = ", conn.ConnId, "[Writer is exit], remote addr is,", conn.GetRemoteAddr().String())
	for {
		select {
		case data := <-conn.MsgChan:
			//have data to send to client
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-conn.ExitChan:
			//reader have exit so writer need exit
			return
		}
	}
}

//SendMsg :pack the msg which will be sent to client
func (conn *Connection) SendMsg(msgId uint32, data []byte) error {
	if conn.ConnStatus == true {
		return errors.New("Connection closed when send msg")
	}
	//pack data | style : MsgDataLen | MsgId | MsgData
	dp := NewDataPack()

	binMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Package err msg ID = ", msgId)
		return errors.New("Pack Error msg")
	}
	//msgData -> client
	conn.MsgChan <- binMsg

	//if _, err := conn.Conn.Write(binMsg); err != nil {
	//	fmt.Println("Write msg id = ", msgId, " error: ", err)
	//	return errors.New("conn Write error")
	//}
	return nil
}

//实现方法 ----------------------------------------------------------------
func (conn *Connection) Start() {
	fmt.Println("conn started ... connId:", conn.ConnId)
	//启动当前链接的读数据业务
	go conn.StartReader()
	//TODO: 启动当前链接的写数据业务
	go conn.StartWriter()
}
func (conn *Connection) Stop() {
	fmt.Println("conn stop() .. ConnId: ", conn.ConnId)
	if conn.ConnStatus == true {
		return
	}
	conn.ConnStatus = true
	//关闭链接
	conn.Conn.Close()
	//send Writer exit
	conn.ExitChan <- true
	//关闭管道
	close(conn.ExitChan)
	close(conn.MsgChan)
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
