package main

import (
	"fmt"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		println("call back before ping error")
	}
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	fmt.Printf("recv from client:data = %s\n", request.GetData())
	if err != nil {
		println("call back ping ping ping error")
	}
}

// Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		println("call back after ping error")
	}
}

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("[zinx V0.3]")
	// 给当前 zinx 框架添加一个自定义的 router
	s.AddRouter(&PingRouter{})
	// 运行服务器
	s.Serve()
}
