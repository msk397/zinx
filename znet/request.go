package znet

import (
	ziface2 "zinx/ziface"
)

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface2.IConnection
	// 客户端请求的数据
	msg ziface2.IMessage
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface2.IConnection {
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
