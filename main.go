package main

import (
	"fmt"
	"github.com/lihuicms-code-rep/zinx/ziface"
	"github.com/lihuicms-code-rep/zinx/znet"
	"mmoGameZinx/core"
)

func main() {
	//创建zinx server 句柄
	var s ziface.IServer
	s = znet.NewServer("MMO Game Server")

	//注册连接创建/销毁的HOOK函数
    s.SetOnConnStart(OnConnectionAdd)
	//注册一些路由业务
	//启动服务
	s.Serve()
}

//客户端建立连接之后的Hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个player对象
	player := core.NewPlayer(conn)

	//给客户端发送数据
	player.SyncPid()
	player.BroadCastBornPosition()
	fmt.Println("====> player pid=", player.Pid, " online.....")
}