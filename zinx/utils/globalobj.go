package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切有关zinx框架的全局参数，供其他模块使用
//一些参数可以通过zinx.json由用户进行配置

type GlobalObj struct {
	TcpServer        ziface.IServer //当前zinx全局的server对象
	Host             string         //当前服务器主机监听的ip
	TcpPort          int            //当前服务器监听的端口
	Name             string         //当前服务器的名字
	Version          string         //当前zinx的版本
	MaxConn          int            //当前服务器的最大连接数
	MaxPacketSize    uint32         //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32         //当前工作业务的Worker池的goroutine数量
	MaxWorkerTaskLen uint32         //当前工作业务的Worker池的goroutine任务队列的最大值
}

//GlobalObject 定义一个全局对外对象
var GlobalObject *GlobalObj

//Reload 从配置文件加载全局对象
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("E:/golang/zinx/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供一个init方法初始化全局对象
func init() {
	//如果配置文件没有加载。默认创建
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		MaxConn:          1000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	//从配置文件加载全局对象
	GlobalObject.Reload()
}
