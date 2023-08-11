package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/ziface"
)

/*
消息处理模块的实现
*/
type MyHandle struct {
	// 存放每个msgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 创建MsgHandle的方法
func NewMsgHandle() *MyHandle {
	return &MyHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度对应的Router处理方法
func (m *MyHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not FOUND! Need Register!")
		return
	}

	//2. 根据MsgID调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 添加具体的处理逻辑
func (m *MyHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1, 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		//id 已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	// 2. 添加msg和API的绑定关系
	m.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %d\n", msgID)
}
