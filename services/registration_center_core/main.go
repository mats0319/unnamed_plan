package main

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/registration_center_core/rpc"
    "github.com/mats9693/unnamed_plan/services/shared/config"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/init"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "net"
)

func main() {
    err := initialize.InitFromConfigCenter(mconst.UID_Service_Registration_center)
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

    err = initAndSetTarget(internetAddress)
    if err != nil {
        mlog.Logger().Error("init or set target failed", zap.Error(err))
        return
    }

    startRegistrationCenter(localAddress)
}

func initAndSetTarget(address string) error {
    client := rpc.GetConfigCenterRCClient()

    err := client.Init(mconfig.GetConfigCenterTarget())
    if err != nil {
        mlog.Logger().Error("init CC-RC client failed", zap.Error(err))
        return err
    }

    err = client.SetRCCoreTarget(address)
    if err != nil {
        mlog.Logger().Error("set rc core target failed", zap.Error(err))
        return err
    }

    return nil
}

func startRegistrationCenter(address string) {
    listener, err := net.Listen("tcp", address)
    if err != nil {
        mlog.Logger().Error(fmt.Sprintf("listen on %s failed", address), zap.Error(err))
        return
    }

    server := grpc.NewServer()
    rpc_impl.RegisterIRegistrationCenterCoreServer(server, rpc.GetRegistrationCenterServer())

    mlog.Logger().Info("> Listening at : " + address)

    // Serve is blocked
    err = server.Serve(listener)
    if err != nil {
        mlog.Logger().Error(fmt.Sprintf("serve on %s failed", address), zap.Error(err))
        return
    }
}
