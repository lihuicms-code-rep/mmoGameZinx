package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lihuicms-code-rep/zinx/ziface"
	"math/rand"
	"mmoGameZinx/pb"
	"sync"
)

//玩家信息
type Player struct {
	Pid  int32
	Conn ziface.IConnection//和客户端的连接
	Pos  Position           //坐标
}

type Position struct {
	X float32 //x轴
	Y float32 //y轴
	Z float32 //z轴
	V float32 //旋转角度
}

var PidGen int32 = 1      //玩家ID计数器,这里简单运用
var PidLock sync.Mutex    //锁

//创建一个玩家
func NewPlayer(conn ziface.IConnection) *Player {
	//生成玩家ID
    PidLock.Lock()
    pid := PidGen
    PidGen++
    defer PidLock.Unlock()

    p := &Player{
    	Pid:pid,
    	Conn:conn,
    	Pos:Position{
    		X:float32(160 + rand.Intn(10)),
    		Y:0,
    		Z:float32(100 + rand.Intn(20)),
    		V:float32(0),
		},
	}

    return p
}


//发送消息给客户端,主要是将proto数据序列化后再调用框架的SendMsg
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err ", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("conn in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("player send msg error ", err)
		return
	}
}


//告知客户端玩家Pid
func (p *Player) SyncPid() {
	//构建msgID=1的消息
	data := &pb.SyncPid{
		Pid:p.Pid,
	}
	//发送消息
	p.SendMsg(1, data)
}

//广播玩家自己的出生位置
func (p *Player) BroadCastBornPosition() {
	data := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,
		Data:&pb.BroadCast_P{
			&pb.Position{
				X:p.Pos.X,
				Y:p.Pos.Y,
				Z:p.Pos.Z,
				V:p.Pos.V,
			},
		},
	}

	p.SendMsg(200, data)
}