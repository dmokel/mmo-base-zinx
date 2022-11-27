package core

import "fmt"

const (
	AOI_MIN_X  int = 80
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 70
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*AOIManager ...
 *the aoi map managerment
 *must CntsX>=CntsY
 */
type AOIManager struct {
	MinX  int
	MaxX  int
	CntsX int
	MinY  int
	MaxY  int
	CntsY int
	grids map[int]*Grid
}

// calculate the width of the grid in the x-axis direction
func (am *AOIManager) gridWidth() int {
	return (am.MaxX - am.MinX) / am.CntsX
}

// calculate the length of the grid in the y-axis direction
func (am *AOIManager) gridLength() int {
	return (am.MaxY - am.MinY) / am.CntsY
}

// NewAOI is the init function for a AOI map
func NewAOI(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	am := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	width := am.gridWidth()
	length := am.gridLength()
	for y := 0; y < am.CntsY; y++ {
		for x := 0; x < am.CntsX; x++ {
			gridID := x + y*am.CntsX
			am.grids[gridID] = NewGrid(gridID,
				am.MinX+x*width,
				am.MinX+(x+1)*width,
				am.MinY+y*length,
				am.MinY+(y+1)*length,
			)
		}
	}

	return am
}

func (am *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager minX: %d, maxX: %d, cntxX:%d, minY: %d, maxY: %d, cntxY: %d, gridIDs:\n",
		am.MinX, am.MaxX, am.CntsX, am.MinY, am.MaxY, am.CntsY)

	for _, grid := range am.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundGridsByGid is the function used to get the surround grids
// 要确保am.CntsX>=am.CntsY，否则会发生错误
func (am *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	if _, ok := am.grids[gid]; !ok {
		return
	}

	grids = append(grids, am.grids[gid])

	idx := gid % am.CntsX

	if idx > 0 {
		grids = append(grids, am.grids[gid-1])
	}

	if idx < am.CntsX-1 {
		grids = append(grids, am.grids[gid+1])
	}

	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, v := range gidsX {
		idy := v / am.CntsY
		if idy > 0 {
			grids = append(grids, am.grids[v-am.CntsX])
		}

		if idy < am.CntsY-1 {
			grids = append(grids, am.grids[v+am.CntsX])
		}
	}

	return
}

// GetGidByPos is the function used to
func (am *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - am.MinX) / am.gridWidth()
	idy := (int(x) - am.MinY) / am.gridLength()

	return idx + idy*am.CntsX
}

// GetPidsByPos is the function used to get surrounding pids around the position
func (am *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	gid := am.GetGidByPos(x, y)

	grids := am.GetSurroundGridsByGid(gid)

	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}

	return
}

// AddPidToGrid means add a player to a grid
func (am *AOIManager) AddPidToGrid(pID, gID int) {
	am.grids[gID].AddPlayer(pID)
}

// RemovePidFromGrid means remove a player from a grid
func (am *AOIManager) RemovePidFromGrid(pID, gID int) {
	am.grids[gID].DeletePlayer(pID)
}

// GetPidsByGid means get all players from a grid
func (am *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = am.grids[gID].GetPlayerIDs()
	return
}

// AddPlayerToGridByPos ...
func (am *AOIManager) AddPlayerToGridByPos(pID int, x, y float32) {
	gID := am.GetGidByPos(x, y)
	grid := am.grids[gID]
	grid.AddPlayer(pID)
}

// RemovePlayerFromGridByPos ...
func (am *AOIManager) RemovePlayerFromGridByPos(pID int, x, y float32) {
	gID := am.GetGidByPos(x, y)
	grid := am.grids[gID]
	grid.DeletePlayer(pID)
}
