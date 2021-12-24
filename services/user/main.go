package main

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/user/config"
	"github.com/mats9693/unnamed_plan/services/user/rpc"
	"github.com/mats9693/utils/toy_server/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", config.GetConfig().Address)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("listen on %s failed", config.GetConfig().Address), zap.Error(err))
		return
	}

	server := grpc.NewServer()
	rpc_impl.RegisterIUserServer(server, rpc.GetUserServer())

	mlog.Logger().Info("> Listening at : " + config.GetConfig().Address)

	// Serve is blocked
	err = server.Serve(listener)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("serve on %s failed", config.GetConfig().Address), zap.Error(err))
		return
	}
}
