package main

import "zinx/znet"

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("[zinx V0.1]")
	// 运行服务器
	s.Serve()
}
