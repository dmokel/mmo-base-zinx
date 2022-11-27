package core

import (
	"log"
	"math/rand"
	"sync"

	"github.com/aceld/zinx/ziface"
	"github.com/dmokel/game-base-zinx/pb"
	"google.golang.org/protobuf/proto"
)

// Player ...
type Player struct {
	Pid  int32
	Conn ziface.IConnection
	X    float32
	Y    float32
	Z    float32
	V    float32
}

var pidGen int32 = 1
var pidLock sync.Mutex

// NewPlayer used to create a new player structure
func NewPlayer(conn ziface.IConnection) *Player {
	pidLock.Lock()
	pid := pidGen
	pidGen++
	pidLock.Unlock()

	return &Player{
		Pid:  pid,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
}

// SendMsg used to send a proto message to the client
func (p *Player) SendMsg(msgID uint32, msg proto.Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("failed to marshal msg, err:", err)
		return
	}

	if p.Conn == nil {
		log.Fatalln("failed to send msg, err: player.conn is nil")
		return
	}

	err = p.Conn.SendMsg(msgID, data)
	if err != nil {
		log.Fatalln("failed to send msg, err:", err)
		return
	}

	return
}

// SyncPid used to send the pid to client
func (p *Player) SyncPid() {
	// construct the msgID:1 message
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// send the message to client
	p.SendMsg(1, protoMsg)
}

//BroadcastPos used to broadcast the player postion
func (p *Player) BroadcastPos() {
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, protoMsg)
}

// Talk used to broadcast the chat content to all players in the game world
func (p *Player) Talk(content string) {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	players := WorldMgr.GetAllPlayers()

	for _, p := range players {
		p.SendMsg(200, msg)
	}
}

// SyncSurrounding used to broadcast the player position to all game players
func (p *Player) SyncSurrounding() {
	pids := WorldMgr.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))

	// broadcast the player position to all other players
	positionMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, pid := range pids {
		players = append(players, WorldMgr.GetPlayerByPid(int32(pid)))
	}
	for _, p := range players {
		p.SendMsg(200, positionMsg)
	}

	// sync all the other players position to the player
	playerMsg := make([]*pb.Player, 0, len(players))
	for _, p := range players {
		playerMsg = append(playerMsg, &pb.Player{
			Pid: p.Pid,
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		})
	}
	syncPlayersMsg := &pb.SyncPlayers{
		Ps: playerMsg,
	}
	p.SendMsg(202, syncPlayersMsg)
}

// UpdatePosition used to update the player position and broadcast the position to all other players
func (p *Player) UpdatePosition(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	players := p.GetSurroundingPlayers()
	for _, p := range players {
		p.SendMsg(200, protoMsg)
	}
}

// GetSurroundingPlayers used to get surrounding players around the player
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgr.AoiMgr.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgr.GetPlayerByPid(int32(pid)))
	}

	return players
}

// Offline ...
func (p *Player) Offline() {
	players := p.GetSurroundingPlayers()
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _, p := range players {
		p.SendMsg(201, protoMsg)
	}

	WorldMgr.AoiMgr.RemovePlayerFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgr.RemovePlayer(p.Pid)
}
