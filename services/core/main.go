package main

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/core/config"
	"github.com/mats9693/unnamed_plan/services/core/db"
	i "github.com/mats9693/unnamed_plan/services/core/init"
	"github.com/mats9693/unnamed_plan/services/core/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"sync"
)

func main() {
	err := initialize.InitFromFile("config.json", mdb.Init, db.Init, config.Init, i.InitSupportedConfig)
	if err != nil {
		mlog.Logger().Error("init failed", zap.Error(err))
		return
	}

	exitChan := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		rpc.GetRegistrationCenterServer().MaintainTarget(exitChan)
		wg.Done()
	}()

	go func() {
		startServer()
		exitChan <- struct{}{}
		wg.Done()
	}()

	wg.Wait()
}

func startServer() {
	address := config.GetConfig().Address

	listener, err := net.Listen("tcp", address)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("listen on %s failed", address), zap.Error(err))
		return
	}

	server := grpc.NewServer()
	rpc_impl.RegisterIConfigCenterServer(server, rpc.GetConfigCenterServer())
	rpc_impl.RegisterIRegistrationCenterCoreServer(server, rpc.GetRegistrationCenterServer())
	//reflection.Register(server) // for grpc ui

	mlog.Logger().Info("> Listening at : " + address)

	// Serve is blocked
	err = server.Serve(listener)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("serve on %s failed", address), zap.Error(err))
		return
	}
}
