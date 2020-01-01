package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lihuicms-code-rep/zinx/ziface"
	"github.com/lihuicms-code-rep/zinx/znet"
	"mmoGameZinx/core"
	"mmoGameZinx/pb"
)

//玩家移动
type MoveAPI struct {
	znet.BaseRouter
}


//具体业务
func (m *MoveAPI) Handle(request ziface.IRequest) {
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Unmarshal msg error", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get pid error", err)
		return
	}

    //给其他玩家广播位置信息,并更新位置坐标
    player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
    player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
