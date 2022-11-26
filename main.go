package main

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/dmokel/game-base-zinx/core"
)

// OnConnAdd ...
func OnConnAdd(conn ziface.IConnection) {
	// create a player
	player := core.NewPlayer(conn)

	// send the message with msgID:1 to client
	player.SyncPid()

	// send the message with msgID:200 to client
	player.BroadcastPos()

	fmt.Printf("player %d is arrived\n", player.Pid)
}

func main() {
	srv := znet.NewServer()

	srv.SetOnConnStart(OnConnAdd)

	srv.Serve()
}
