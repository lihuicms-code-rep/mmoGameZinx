package core

import (
	"fmt"
	"sync"
)

//AOI格子
type Grid struct {
	GID       int          //格子ID
	MinX      int          //格子左上
	MaxX      int          //右上
	MinY      int          //左下
	MaxY      int          //右下
	playerIDs map[int]bool //格子内玩家
	pIDLock   sync.RWMutex //读写锁
}


//创建一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:gID,
		MinX:minX,
		MinY:minY,
		MaxX:maxX,
		MaxY:maxY,
		playerIDs:make(map[int]bool),

	}
}


//给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

//从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

//得到当前格子中所有玩家
func (g *Grid) GetPlayerIDs() []int {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	var playerIDs []int
	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return playerIDs
}


//调试:打印格子基本信息
//重写String方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		                        g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
