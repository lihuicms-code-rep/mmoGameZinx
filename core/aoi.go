package core

import "fmt"

//AOI区域管理模块

type AOIManager struct {
	MinX  int           //区域左边坐标
	MaxX  int           //区域右边坐标
	CntsX int           //X方向格子数量
	MinY  int           //区域上边坐标
	MaxY  int           //区域下边坐标
	CntsY int           //Y方向格子数量
	grids map[int]*Grid //区域中有哪些格子
}

//初始化
func NewAOIManager(minX, maxX, minY, maxY, cntsX, cntsY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//给AOI区域中格子进行编号并初始化
	//具体的格子编号可以参考图像进行归纳
	for y := 0; y < cntsY; y++ { //从上至下
		for x := 0; x < cntsX; x++ { //从左至右
			//格子编号: id = idy *cntsX + idx
			gid := y*cntsX + x
			aoiManager.grids[gid] = NewGrid(gid,
				aoiManager.MinX+x*aoiManager.gridWidth(),
				aoiManager.MinX+(x+1)*aoiManager.gridWidth(),
				aoiManager.MinY+y*aoiManager.gridLength(),
				aoiManager.MinY+(y+1)*aoiManager.gridLength())
		}
	}

	return aoiManager
}

//得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个给子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//调试信息,打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager Info: MinX:%d, MaxX:%d, MinY:%d, MaxY:%d, cntsX:%d, cntsY:%d\n",
		m.MinX, m.MaxX, m.MinY, m.MaxY, m.CntsX, m.CntsY)

	for _, grid := range m.grids {
		s += fmt.Sprintf("%s\n", grid)
	}

	return s
}

//根据格子ID,得到周边九宫格格子集合
//计算过程以格子编号为8带入体会
func (m *AOIManager) GetSurroundGridsByGid(gID int) []*Grid {
	//判断当前gID是否在该区域中
	if _, ok := m.grids[gID]; !ok {
		return nil
	}

	//初始化返回切片
	grids := make([]*Grid, 0)
	grids = append(grids, m.grids[gID]) //8

	//判断gID左边是否有格子？右边是否有格子？
	idx := gID % m.CntsX
	if idx > 0 {
		grids = append(grids, m.grids[gID-1]) //7
	}

	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1]) //9
	}

	//将X轴当前格子都取出
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID) // 7,8,9
	}

	//判断每个格子的上下是否有格子
	for _, v := range gridsX {
		idy := v / m.CntsY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}

		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}

	return grids
}

//通过x,y坐标得到周边九宫格内所有playerID
func (m *AOIManager) GetPidsByPos(x, y float32) []int {
	//得到当前玩家所属的格子ID
	gID := m.GetGidByPos(x, y)

	//根据格子ID得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gID)

	//获取九宫格信息中全部playerID信息
	playerIDs := make([]int, 0)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}

	return playerIDs
}

//通过x,y坐标得到当前所属格子ID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

//添加一个playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

//移除一个格子中的playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

//通过GID获取全部playerID
func (m *AOIManager) GetPidsByGid(gID int) []int {
	return m.grids[gID].GetPlayerIDs()
}

//通过坐标把player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Add(pID)
}

//通过坐标把一个player从格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Remove(pID)
}
