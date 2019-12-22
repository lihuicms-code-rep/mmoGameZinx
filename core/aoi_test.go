package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化
	aoiMgr := NewAOIManager(0, 250, 0, 250, 5, 5)

	//调试信息
	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 0, 250, 5, 5)

	for gid, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Printf("当前格子ID=%d, 它的周围格子有%d个", gid, len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}

		fmt.Println("surrounding grid ID ", gIDs)
	}
}
