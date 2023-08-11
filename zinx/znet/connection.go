package znet

import (
	"fmt"
	"net"
	"zinx/zinx/ziface"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true
	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		return
	}

	// 通知从缓冲队列读数据的业务，该链接已经关闭
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer func() {
		c.Stop()
		fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	}()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//显示客户端发送的内容
		//fmt.Printf("recv client buf %s, cnt = %d\n", buf, cnt)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//执行注册的路由方法
		go func(req ziface.IRequest) {
			//从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)

	}
}

func (c *Connection) StartWriter() {

}
