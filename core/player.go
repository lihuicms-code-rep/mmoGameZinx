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
	Conn ziface.IConnection //和客户端的连接
	Pos  Position           //坐标
}

type Position struct {
	X float32 //x轴
	Y float32 //y轴
	Z float32 //z轴
	V float32 //旋转角度
}

var PidGen int32 = 1   //玩家ID计数器,这里简单运用
var PidLock sync.Mutex //锁

//创建一个玩家
func NewPlayer(conn ziface.IConnection) *Player {
	//生成玩家ID
	PidLock.Lock()
	pid := PidGen
	PidGen++
	defer PidLock.Unlock()

	p := &Player{
		Pid:  pid,
		Conn: conn,
		Pos: Position{
			X: float32(160 + rand.Intn(10)),
			Y: 0,
			Z: float32(100 + rand.Intn(20)),
			V: float32(0),
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
		Pid: p.Pid,
	}
	//发送消息
	p.SendMsg(1, data)
}

//广播玩家自己的出生位置
func (p *Player) BroadCastBornPosition() {
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.Pos.X,
				Y: p.Pos.Y,
				Z: p.Pos.Z,
				V: p.Pos.V,
			},
		},
	}

	p.SendMsg(200, data)
}

//世界广播聊天信息
func (p *Player) Talk(content string) {
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//得到当前世界所有玩家
	players := WorldMgrObj.GetAllPlayers()

	//发送消息
	for _, p := range players {
		p.SendMsg(2, data)
	}
}

//同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	//获取当前玩家周围玩家(九宫格)
	players := p.GetSurroundPlayers()

	//给周围玩家发送自己的位置信息
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.Pos.X,
				Y: p.Pos.Y,
				Z: p.Pos.Z,
				V: p.Pos.V,
			},
		},
	}

	for _, player := range players {
		player.SendMsg(200, data)
	}

	//将周围玩家的信息发送给自己
	playerMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		playerMsg = append(playerMsg, &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.Pos.X,
				Y: player.Pos.Y,
				Z: player.Pos.Z,
				V: player.Pos.V,
			},
		})
	}

	playersMsg := &pb.SyncPlayers{
		Ps:playerMsg,
	}

	p.SendMsg(202, playersMsg)
}


//广播当前玩家位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	//更新新的位置
	p.Pos.X = x
	p.Pos.Y = y
	p.Pos.Z = z
	p.Pos.V = v

	//广播协议
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.Pos.X,
				Y: p.Pos.Y,
				Z: p.Pos.Z,
				V: p.Pos.V,
			},
		},
	}

	//获取周边全部玩家(九宫格)
    players := p.GetSurroundPlayers()
    for _, player := range players {
    	player.SendMsg(200, data)
	}
}


//获取玩家九宫格内所有玩家信息
func (p *Player) GetSurroundPlayers() []*Player {
	pids := WorldMgrObj.AOIMgr.GetPidsByPos(p.Pos.X, p.Pos.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//玩家下线
func (p *Player) Offline() {
	players := p.GetSurroundPlayers()
	data := &pb.SyncPid{
		Pid:p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, data)
	}

	//同时将当前玩家从AOI管理信息中删除
	WorldMgrObj.RemovePlayerByPid(p.Pid)
	WorldMgrObj.AOIMgr.RemoveFromGridByPos(int(p.Pid), p.Pos.X, p.Pos.Z)
}