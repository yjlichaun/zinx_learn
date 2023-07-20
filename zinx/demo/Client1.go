package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("Client1 start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error: ", err)
		return
	}
	for {
		//send pack message
		dp := znet.NewDataPack()
		binary, err := dp.Pack(znet.NewMsgPackage(1, []byte("ZinxV0.6 client1 Test Message")))
		if err != nil {
			fmt.Println("pack error: ", err)
			return
		}
		_, err = conn.Write(binary)
		if err != nil {
			fmt.Println("write error: ", err)
			return
		}
		//server return a messageData
		// read Head get ID and DataLen
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("conn read error: ", err)
			return
		}
		msgData, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Client unpack error: ", err)
			return
		}
		if msgData.GetMsgLen() > 0 {
			msg := msgData.(*znet.Message)
			msg.Data = make([]byte, msgData.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("conn read error: ", err)
				return
			}
			fmt.Println("————> Recv Server Message ID:", msg.GetMessageId(), " MsgLen:", msg.GetMsgLen(), " Data:", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}

}
