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

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写 ping...ping...ping
	fmt.Printf("recv from client, msgId: %d, data: %s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("[zinx V0.5]")
	// 给当前 zinx 框架添加一个自定义的 router
	s.AddRouter(&PingRouter{})
	// 运行服务器
	s.Serve()
}
