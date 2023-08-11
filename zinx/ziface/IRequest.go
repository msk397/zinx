package ziface

/*
IRequest 接口：
实际上是把客户端请求的连接信息和请求的数据包装到了 Request 里
*/
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgID 得到请求的消息 ID
	GetMsgID() uint32

	// GetMsgLen 得到请求的消息长度
	GetMsgLen() uint32
}
