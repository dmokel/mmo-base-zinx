package core

import (
	"fmt"
	"testing"
)

// TestNewAOIManager test the aoi-manager init function logic
func TestNewAOIManager(t *testing.T) {
	am := NewAOI(100, 300, 4, 200, 450, 5)

	fmt.Println(am)
}

func TestGetSurroundGridsByGid(t *testing.T) {
	am := NewAOI(100, 350, 5, 200, 450, 5)

	for gid := range am.grids {
		grids := am.GetSurroundGridsByGid(gid)
		gids := make([]int, 0, len(grids))
		for _, v := range grids {
			gids = append(gids, v.GID)
		}

		fmt.Printf("gid:%d, grids num:%d, gids are:%v\n", gid, len(grids), gids)
	}
}
