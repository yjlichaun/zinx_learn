package main

import "zinx/znet"

func main() {
	z := znet.NewServer("[demoServer]")
	z.Serve()
}
