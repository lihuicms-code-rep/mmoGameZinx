package main

import "github.com/lihuicms-code-rep/zinx/znet"

func main() {
	//创建zinx server 句柄
	var s znet.Server
	s = znet.NewServer("MMO Game Server")

	//注册连接创建/销毁的HOOK函数
	//注册一些路由业务
	//启动服务
	s.Serve()
}