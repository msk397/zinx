package test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/zinx/znet"
)

// 测试datapack拆包和封包
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	//1. 创建socketTCp
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// go 负责处理从客户端处理业务
	//2. 从客户端读取数据，拆包处理
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				/* 拆包的过程 */
				//定义一个拆包对象dp
				dp := znet.NewDataPack()
				for {
					// 1. 第一次从conn读head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						return
					}

					message, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					if message.GetMsgLen() > 0 {
						//说明msg里面有数据，需要根据msg里面的msglen接着读

						//从conn里面接着读，根据msg里面的msglen
						msg := message.(*znet.Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						fmt.Printf("--> Recv MsgID: %d, MsgLen:%d, data:%s\n", msg.GetMsgId(), msg.GetMsgLen(), msg.GetData())

					}

				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	//创建一个封包对象 dp
	dp := znet.NewDataPack()
	//模拟粘包过程，封装两个msg一起发送
	//封第一个
	msg1 := &znet.Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	pack1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	//封第二个
	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!', '!'},
	}
	pack2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
		return
	}
	//粘起来
	sendData := append(pack1, pack2...)
	//发送
	conn.Write(sendData)

	//阻塞
	select {}
}
