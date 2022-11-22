package core

import (
	"fmt"
	"sync"
)

/*Grid ...
 *thie grid of AOI map
 */
type Grid struct {
	// grid id
	GID int
	// left boundary coordinate
	MinX int
	// right boundary coordinate
	MaxX int
	// lower boundary coordinate
	MinY int
	// upper boundary coordinate
	MaxY int
	// the collection of current players
	playerIDs map[int]bool
	pIDLock   sync.RWMutex
}

// NewGrid is the init function for the grid
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// AddPlayer is the function used to add a player into the grid
func (g *Grid) AddPlayer(pID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[pID] = true
}

// DeletePlayer is the function used to delete a player from the grid
func (g *Grid) DeletePlayer(pID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, pID)
}

// GetPlayerIDs is the function used to return all the player's IDs in the grid
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for pID := range g.playerIDs {
		playerIDs = append(playerIDs, pID)
	}
	return
}

// for debug
func (g *Grid) String() string {
	return fmt.Sprintf("Grid ID: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v\n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
