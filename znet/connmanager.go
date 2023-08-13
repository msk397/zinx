package znet

import (
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection

	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (c *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将conn加入到ConnManager中
	c.connections[conn.GetConnID()] = conn

	// 打印连接总数
	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())
}

// Remove 删除链接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除连接信息
	delete(c.connections, conn.GetConnID())

	// 打印连接总数
	fmt.Println("connection Remove ConnID = ", conn.GetConnID(), " successfully: conn num = ", c.Len())
}

// Get 根据connID获取链接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	// 根据connID获取连接信息
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, fmt.Errorf("connection not FOUND")
	}
}

// Len 得到当前链接总数
func (c *ConnManager) Len() uint {
	return uint(len(c.connections))
}

// ClearConn 清除并终止所有链接
func (c *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range c.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(c.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}
