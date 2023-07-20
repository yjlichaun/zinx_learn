package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

//构造方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//Pack Method
//style:|dataLen|msgID|data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//create a buf to hold bytes
	dataBuff := bytes.NewBuffer([]byte{})
	// dataLen -> dataBuff
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}
	// MsgId -> dataBuff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMessageId())
	if err != nil {
		return nil, err
	}
	// data -> dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//Unpack Method
//get packageHead message , head -> dataLen to read
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//create a io.Reader to input binary data
	dataBuff := bytes.NewReader(binaryData)
	//only decompression head message ,get dataLen and MsgId
	msg := &Message{}
	//read dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	//read MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//judge dataLen is or not over the maxPacket length we limit
	if utils.GlobalObject.MaxPacketSize > 0 && msg.MsgLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("dataLen is over the maxPacket length we limit")
	}

	return msg, nil

}
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4bytes) + ID uint32(4bytes)
	return 8
}
