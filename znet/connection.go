package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前Connection隶属于哪个Server
	TcpServer ziface.IServer

	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool

	// 无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte

	//该链接处理的方法Router
	MsgHandler ziface.IMsgHandle
}

// NewConnection 初始化连接模块的方法
func NewConnection(s ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  s,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandle,
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()
	// 按照开发者传递进来的创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error msg id = %d, err = %s\n", msgId, err)
		return err
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 调用开发者注册的连接断开之后需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		return
	}

	//告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 通知从缓冲队列读数据的业务，该链接已经关闭
	close(c.ExitChan)
	close(c.msgChan)

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

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer func() {
		c.Stop()
		fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	}()

	for {
		//读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData) //ReadFull 会把headData填充满为止
		if err != nil {
			fmt.Println("read head error", err)
			break
		}
		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据datalen 再次读取data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//显示客户端发送的内容
		//fmt.Printf("recv client buf %s, cnt = %d\n", buf, cnt)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//说明已经开启了工作池机制，发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//执行注册的路由方法
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

/*
写消息Goroutine，专门发送给客户端消息的模块
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("Send data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}
