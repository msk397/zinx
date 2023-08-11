package main

import (
	"fmt"
	"net"
	"time"
	"zinx/zinx/znet"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(1 * time.Second)
	// 创建一个客户端句柄
	dial, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for i := 0; i < 1000; i++ {
		// 发送封包的 msg 消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(uint32(i%3), []byte("ZinxV0.6 client Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := dial.Write(binaryMsg); err != nil {
			fmt.Println("Write error:", err)
			return
		}
		// 服务器就应该给我们回复一个 message 数据，msgID:1 ping...ping...ping
		// 1. 先读取流中的 head 部分，得到 ID 和 dataLen
		headData := make([]byte, dp.GetHeadLen())
		if _, err := dial.Read(headData); err != nil {
			fmt.Println("Read head error:", err)
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Server unpack err:", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 2. 再根据 dataLen 进行第二次读取，将 data 读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := dial.Read(msg.Data); err != nil {
				fmt.Println("Server unpack data err:", err)
				return
			}
			fmt.Printf("Recv from server MsgID: %d, data: %s\n", msg.Id, msg.Data)
		}
		time.Sleep(1 * time.Second)
	}
}
