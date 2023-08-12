package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	ziface "zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个msgID对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负责worker取任务的消息队列
	WorkerPoolSize uint32

	//业务工作WorkerPool的worker数量
	TaskQueue []chan ziface.IRequest
}

// NewMsgHandle 创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度对应的Router处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
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
func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1, 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		//id 已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	// 2. 添加msg和API的绑定关系
	m.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %d\n", msgID)
}

// 启动一个Worker工作池
func (m *MsgHandle) StartWorkerPool() {
	//根据WorkerPoolSize分别开启Worker， 每个Worker用一个go承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//启动worker
		// 1. 当前worker对应的channel消息队列， 开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker， 阻塞等待消息从channel传递
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ", workerID, " is Started...")

	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker ,根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(), "to workerID = ", workerID)
	// 将消息发送给对应worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
