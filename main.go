package main

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/dmokel/game-base-zinx/api"
	"github.com/dmokel/game-base-zinx/core"
)

// OnConnAdd ...
func OnConnAdd(conn ziface.IConnection) {
	// create a player
	player := core.NewPlayer(conn)

	// send the message with msgID:1(palyer id) to client
	player.SyncPid()

	// send the message with msgID:200(player postion) to client
	player.BroadcastPos()

	// sync postion with the surrounding players
	player.SyncSurrounding()

	// add the palyer to world
	core.WorldMgr.AddPlayer(player)

	// bind the pid to conn
	conn.SetProperty("pid", player.Pid)

	fmt.Printf("player %d is arrived\n", player.Pid)
}

// OnConnLost ...
func OnConnLost(conn ziface.IConnection) {
	pid, err := conn.GetProperty("pid")
	if err != nil {
	}

	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.Offline()
}

func main() {
	srv := znet.NewServer()
	srv.SetOnConnStart(OnConnAdd)

	srv.AddRouter(2, &api.WorldChat{})
	srv.AddRouter(3, &api.Move{})

	srv.Serve()
}
