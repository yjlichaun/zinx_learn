package znet

type Message struct {
	Id     uint32 //消息ID
	MsgLen uint32 //消息数据长度
	Data   []byte //消息内容
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:     id,
		MsgLen: uint32(len(data)),
		Data:   data,
	}
}
func (m *Message) GetMessageId() uint32 {
	return m.Id
}
func (m *Message) SetMessageId(id uint32) {
	m.Id = id
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}
func (m *Message) SetMsgLen(len uint32) {
	m.MsgLen = len
}
