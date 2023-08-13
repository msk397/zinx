package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	ziface "zinx/ziface"
)

// Server 服务器模块
// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int

	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务的API关系
	MsgHandler ziface.IMsgHandle

	// 该server的连接管理器
	ConnMgr ziface.IConnmanager

	// 该Server创建连接之后，自动调用hook函数
	OnConnStart func(conn ziface.IConnection)

	// 该Server断开连接之后，自动调用hook函数
	OnConnStop func(conn ziface.IConnection)
}

// AddRouter 添加路由方法
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[zinx] Server Name: %s, listenner at IP: %s, Port %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	// 初始化消息队列及Worker
	s.MsgHandler.StartWorkerPool()
	go func() {
		//1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr err", err)
			return
		}
		//2. 启动服务器的监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP err", err)
			return
		}
		fmt.Printf("Start Zinx server %s success, now listenning...\n", s.Name)

		var cid uint32
		cid = 0

		//3. 阻塞等待客户端的连接，处理客户 端的连接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//最大连接数判断
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端做一个响应

				fmt.Println("too many Connections")
				err := conn.Close()
				if err != nil {
					return
				}
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx Server ")
	s.ConnMgr.ClearConn()

}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

// GetConnMgr 获取ConnMgr
func (s *Server) GetConnMgr() ziface.IConnmanager {
	return s.ConnMgr
}

// NewServer 初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("==> Call ConnStart....")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("==> Call ConnStop....")
		s.OnConnStop(connection)
	}
}
