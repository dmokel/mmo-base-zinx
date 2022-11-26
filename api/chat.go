package api

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/dmokel/game-base-zinx/core"
	"github.com/dmokel/game-base-zinx/pb"
	"google.golang.org/protobuf/proto"
)

// WorldChat ...
type WorldChat struct {
	znet.BaseRouter
}

// Handle ...
func (wc *WorldChat) Handle(req ziface.IRequest) {
	protoMsg := &pb.Talk{}
	proto.Unmarshal(req.GetData(), protoMsg)

	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
	}

	p := core.WorldMgr.GetPlayerByPid(pid.(int32))

	p.Talk(protoMsg.Content)
}
