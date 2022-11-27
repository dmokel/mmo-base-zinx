package api

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/dmokel/game-base-zinx/core"
	"github.com/dmokel/game-base-zinx/pb"
	"google.golang.org/protobuf/proto"
)

// Move ...
type Move struct {
	znet.BaseRouter
}

// Handle ...
func (m *Move) Handle(req ziface.IRequest) {
	protoMsg := &pb.Position{}
	if err := proto.Unmarshal(req.GetData(), protoMsg); err != nil {
		fmt.Println("Move: failed to Unmarshal Position, err:", err)
		return
	}

	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Move: failed to GetProperty pid, err:", err)
		return
	}

	fmt.Printf("Move: playerId=%d, move to (%f,%f,%f,%f)\n", pid.(int32), protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)

	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.UpdatePosition(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
