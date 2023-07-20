package test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/znet"
)

// only test dataPack pack and unpack methods
func TestDataPack(t *testing.T) {
	// imitate server
	// create socketTcp
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen error", err)
		return
	}
	//create a goroutine to Responsible for handling business from the client
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
				return
			}
			go func(conn net.Conn) {
				//Responsible for handling business from the client
				//-------------------pack-------------------------------
				//create a object to pack
				dp := znet.NewDataPack()
				for {
					//the first read from conn ,get the package head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg is not empty ,need to read second
						//the second read for conn, get data with package head.dataLen
						msg := msgHead.(*znet.Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//with dataLen to read data from ioStream
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error", err)
							return
						}
						//println
						fmt.Println("————> Recv MsgId:", msg.Id, "dataLen = ", msg.MsgLen, "data = ", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	//c client
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial error", err)
		return
	}
	//create a object to unpack
	dp := znet.NewDataPack()
	//imitate Sticky wrapping process , pack two msg send together
	//first pack
	firstMsg := &znet.Message{
		Id:     1,
		MsgLen: 4,
		Data:   []byte{'z', 'i', 'n', 'x'},
	}
	firstSendData, err := dp.Pack(firstMsg)
	if err != nil {
		fmt.Println("client pack firstMsg error", err)
		return
	}
	//second pack
	secondMsg := &znet.Message{
		Id:     2,
		MsgLen: 7,
		Data: []byte{
			'n', 'i', 'h', 'a', 'o', '!', '!',
		},
	}
	secondSendData, err := dp.Pack(secondMsg)
	if err != nil {
		fmt.Println("client pack secondMsg error", err)
		return
	}
	//unite two package
	sendData := append(firstSendData, secondSendData...)
	//send data to server
	_, err = conn.Write(sendData)
	if err != nil {
		fmt.Println("client write error", err)
		return
	}
	//block client
	select {}
}
