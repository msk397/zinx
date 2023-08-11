package znet

import "zinx/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求消息的 ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

// GetMsgLen 获取请求消息的长度
func (r *Request) GetMsgLen() uint32 {
	return r.msg.GetMsgLen()
}
