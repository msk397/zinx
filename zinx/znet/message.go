package znet

type Message struct {
	Id      uint32 // 消息的 ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(u uint32) {
	m.Id = u
}

func (m *Message) SetDataLen(u uint32) {
	m.DataLen = u
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
