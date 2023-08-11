package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(1 * time.Second)
	// 创建一个客户端句柄
	dial, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for i := 0; i < 1000; i++ {
		// 发送数据
		str := fmt.Sprintf("hello zinx v0.1, index: %d", i)
		_, err := dial.Write([]byte(str))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}
		// 接收数据
		buf := make([]byte, 512)
		cnt, err := dial.Read(buf)
		if err != nil {
			fmt.Println("Receive buf err", err)
			return
		}
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
