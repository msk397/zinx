package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	ziface2 "zinx/ziface"
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

	MsgHandler ziface2.IMsgHandle
}

// AddRouter 添加路由方法
func (s *Server) AddRouter(msgID uint32, router ziface2.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[zinx] Server Name: %s, listenner at IP: %s, Port %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
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
			// 将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或者回收

}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface2.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}