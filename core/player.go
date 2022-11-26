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