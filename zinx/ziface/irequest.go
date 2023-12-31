package ziface

//IRequest 接口
//实际上是把客户端请求的链接数据，和请求的数据包装到一个Request中

type IRequest interface {
	//GetConnection 得到当前链接
	GetConnection() IConnection
	//GetData 得到请求的消息数据
	GetData() []byte
	//GetMsgId get the message id
	GetMsgId() uint32
}
