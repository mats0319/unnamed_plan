package main

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/note/config"
	"github.com/mats9693/unnamed_plan/services/note/db"
	"github.com/mats9693/unnamed_plan/services/note/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	initialize.InitFromFile("config.json", mdb.Init, config.Init, db.Init)

	listener, err := net.Listen("tcp", config.GetConfig().Address)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("listen on %s failed", config.GetConfig().Address), zap.Error(err))
		return
	}

	server := grpc.NewServer()
	noteServer, err := rpc.GetNoteServer(config.GetConfig().UserServerAddress)
	if err != nil {
		mlog.Logger().Error("init note server failed", zap.Error(err))
		return
	}
	rpc_impl.RegisterINoteServer(server, noteServer)

	mlog.Logger().Info("> Listening at : " + config.GetConfig().Address)

	// Serve is blocked
	err = server.Serve(listener)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("serve on %s failed", config.GetConfig().Address), zap.Error(err))
		return
	}
}
