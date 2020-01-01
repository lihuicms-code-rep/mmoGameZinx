package core

import "sync"

//当前游戏世界管理模块
type WorldManager struct {
	AOIMgr  *AOIManager       //aoi管理模块
	Players map[int32]*Player //当前在线所有玩家
	pLock   sync.RWMutex      //保护Player的读写锁
}

//全局变量,对外的世界管理模块
var WorldMgrObj *WorldManager

//初始化
func init() {
	WorldMgrObj = &WorldManager{
		AOIMgr:NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_X, AOI_CNTS_Y),
		Players:make(map[int32]*Player),
	}
}

//添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()

	wm.Players[player.Pid] = player

	//同时将玩家添加入AOIMgr
	wm.AOIMgr.AddToGridByPos(int(player.Pid), player.Pos.X, player.Pos.Z)
}


//删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	player := wm.Players[pid]
	wm.pLock.Lock()
	defer wm.pLock.Unlock()

	//将玩家从AOIMgr删除
	wm.AOIMgr.RemoveFromGridByPos(int(player.Pid), player.Pos.X, player.Pos.Z)

	//再从世界管理中删除
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	delete(wm.Players, player.Pid)
}

//通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

//获取全部在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, p := range wm.Players {
		players = append(players, p)
	}

	return players
}