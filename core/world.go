package core

import "sync"

// WorldManager used to manage the game world
type WorldManager struct {
	AoiMgr  *AOIManager
	Players map[int32]*Player
	plock   sync.RWMutex
}

// WorldMgr ...
var WorldMgr *WorldManager

func init() {
	WorldMgr = &WorldManager{
		AoiMgr:  NewAOI(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

// AddPlayer used to add a player to players
func (w *WorldManager) AddPlayer(p *Player) {
	w.plock.Lock()
	w.Players[p.Pid] = p
	w.plock.Unlock()

	w.AoiMgr.AddPlayerToGridByPos(int(p.Pid), p.X, p.Y)
}

// RemovePlayer used to delete a player from players
func (w *WorldManager) RemovePlayer(pid int32) {
	p := w.Players[pid]
	w.AoiMgr.RemovePlayerFromGridByPos(int(p.Pid), p.X, p.Y)

	w.plock.Lock()
	delete(w.Players, pid)
	w.plock.Unlock()
}

// GetPlayerByPid used to get a player by playerId
func (w *WorldManager) GetPlayerByPid(pid int32) *Player {
	w.plock.Lock()
	defer w.plock.Unlock()
	return w.Players[pid]
}

// GetAllPlayers used to get all players in the game world
func (w *WorldManager) GetAllPlayers() []*Player {
	w.plock.Lock()
	defer w.plock.Unlock()

	players := make([]*Player, 0)

	for _, p := range w.Players {
		players = append(players, p)
	}

	return players
}
