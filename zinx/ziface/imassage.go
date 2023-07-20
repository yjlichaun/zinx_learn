package ziface

type IMessage interface {
	//GetMessageId 获取消息id
	GetMessageId() uint32
	//SetMessageId 设置消息id
	SetMessageId(id uint32)
	//GetMsgLen 获取数据长度
	GetMsgLen() uint32
	//SetMsgLen 设置数据长度
	SetMsgLen(len uint32)
	//GetData 获取数据
	GetData() []byte
	//SetData 设置消息内容
	SetData(data []byte)
}
