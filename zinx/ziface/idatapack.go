package ziface

//解决tcp粘包问题
//封包，拆包模块面向tcp链接中的数据流

type IDataPack interface {
	//GetHeadLen 获取包头长度的方法
	GetHeadLen() uint32
	//Pack 封包方法
	Pack(message IMessage) ([]byte, error)
	//Unpack 拆包方法
	Unpack(data []byte) (IMessage, error)
}
