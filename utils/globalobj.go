package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
存储一切有关 Zinx 框架的全局参数，供其他模块使用
一些参数是可以通过 zinx.json 由用户进行配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	TCPServer ziface.IServer // 当前 Zinx 全局的 Server 对象
	Host      string         // 当前服务器主机监听的 IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	/*
		Zinx
	*/
	Version        string // 当前 Zinx 的版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前 Zinx 框架数据包的最大值
}

/*
定义一个全局的对外 GlobalObj
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供一个 init 方法，初始化当前的 GlobalObject
*/
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.6",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPackageSize: 4096,
	}

	// 应该尝试从 conf/zinx.json 去加载一些用户自定义的参数
	GlobalObject.Reload()
}
