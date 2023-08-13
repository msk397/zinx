package main

import (
	"fmt"
	"zinx/ziface"
	znet2 "zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet2.BaseRouter
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

type HelloRouter struct {
	znet2.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Printf("recv from client, msgId: %d, data: %s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("Hello Zinx V0.10"))
	if err != nil {
		fmt.Println(err)
	}
}

type QuitRouter struct {
	znet2.BaseRouter
}

func (q *QuitRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call QuitRouter Handle...")
	fmt.Printf("recv from client, msgId: %d, data: %s\n", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("Quit Zinx V0.10"))
	if err != nil {
		fmt.Println(err)
	}
}

func ConnStart(conn ziface.IConnection) {
	fmt.Println("Conn Start...")
	err := conn.SendMsg(202, []byte("Hello Zinx V0.10"))
	if err != nil {
		fmt.Println(err)
	}

	// 给当前的连接设置一些属性
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "ZinxV0.10")
	conn.SetProperty("Home", "test.com")

}

func ConnStop(conn ziface.IConnection) {
	fmt.Println("Conn Stop...")

	// 获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

}

func main() {
	// 创建一个服务器句柄
	s := znet2.NewServer()

	// 注册连接的 hook 钩子函数
	s.SetOnConnStart(ConnStart)
	s.SetOnConnStop(ConnStop)

	// 给当前 zinx 框架添加一个自定义的 router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.AddRouter(2, &QuitRouter{})

	// 运行服务器
	s.Serve()
}
