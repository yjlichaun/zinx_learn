package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error: ", err)
		return
	}
	for {
		//链接调用write写数据
		_, err := conn.Write([]byte("Hello Zinx V0.1"))
		if err != nil {
			fmt.Println("write conn error: ", err)
			return
		}
		//链接调用read读取数据
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read conn error: ", err)
			return
		}
		fmt.Printf("server call back : %s, cnt = %d\n", buf, n)
		time.Sleep(1 * time.Second)
	}

}
