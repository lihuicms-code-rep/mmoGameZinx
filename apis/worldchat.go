package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lihuicms-code-rep/zinx/ziface"
	"github.com/lihuicms-code-rep/zinx/znet"
	"mmoGameZinx/core"
	"mmoGameZinx/pb"
)

//世界聊天路由业务
type WorldChatAPI struct {
	znet.BaseRouter
}


//具体业务
func (wc *WorldChatAPI) Handle(request ziface.IRequest) {
    //1.解析客户端传递的proto
    msg := &pb.Talk{}
    err := proto.Unmarshal(request.GetData(), msg)
    if err != nil {
    	fmt.Println("Talk msg Unmarshal error")
		return
	}

    //2.当前聊天数据来自于哪个玩家发送的
    pid, err := request.GetConnection().GetProperty("pid")
    if err != nil {
    	fmt.Println("get pid error", err)
	}

    player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

    //3.广播
    player.Talk(msg.GetContent())
}
