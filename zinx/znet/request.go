package znet

import "zinx/ziface"

type Request struct {
	//已经和客户端建立链接的connection
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

func (req *Request) GetConnection() ziface.IConnection {
	return req.conn
}
func (req *Request) GetData() []byte {
	return req.msg.GetData()
}
func (req *Request) GetMsgId() uint32 {
	return req.msg.GetMessageId()
}
