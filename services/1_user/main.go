package main

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/1_user/config"
	"github.com/mats9693/unnamed_plan/services/1_user/db"
	"github.com/mats9693/unnamed_plan/services/1_user/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	err := initialize.InitFromConfigCenter(mconst.UID_Service_User, mdb.Init, config.Init, db.Init)
	if err != nil {
		mlog.Logger().Error("init failed", zap.Error(err))
		return
	}

	ip, port, err := initialize.GetIPAndFreePort()
	if err != nil {
		mlog.Logger().Error("get ip or free port failed", zap.Error(err))
		return
	}

	localAddress := fmt.Sprintf("127.0.0.1:%d", port)
	internetAddress := fmt.Sprintf("%s:%d", ip, port)

	rceServer, err := rce.InitAndRegister(mconfig.GetCoreTarget(), mconst.UID_Service_User, internetAddress)
	if err != nil {
		mlog.Logger().Error("register server to rc failed", zap.Error(err))
		return
	}

	startUserServer(localAddress, rceServer)
}

func startUserServer(address string, rceServer rpc_impl.IRegistrationCenterEmbeddedServer) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("listen on %s failed", address), zap.Error(err))
		return
	}

	server := grpc.NewServer()
	rpc_impl.RegisterIUserServer(server, rpc.GetUserServer())
	rpc_impl.RegisterIRegistrationCenterEmbeddedServer(server, rceServer)
	//reflection.Register(server) // for grpc ui

	mlog.Logger().Info("> Listening at : " + address)

	// Serve is blocked
	err = server.Serve(listener)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("serve on %s failed", address), zap.Error(err))
		return
	}
}
