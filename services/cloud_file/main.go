package main

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/cloud_file/config"
    "github.com/mats9693/unnamed_plan/services/cloud_file/rpc"
    "github.com/mats9693/unnamed_plan/services/shared/config"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "net"
    "os"
)

func init() {
    root := config.GetConfig().CloudFileRootPath
    if len(root) < 1 {
        root = mconfig.GetExecDir()+"files/"
    }
    path := utils.FormatDirSuffix(root) + config.GetConfig().CloudFilePublicDir

    err := os.MkdirAll(path, 0755)
    if err != nil {
        mlog.Logger().Error("os.MkdirAll failed", zap.Error(err))
        os.Exit(-1)
    }

    mlog.Logger().Info("> Cloud file directory init finish.")
}

func main() {
    listener, err := net.Listen("tcp", config.GetConfig().Address)
    if err != nil {
        mlog.Logger().Error(fmt.Sprintf("listen on %s failed", config.GetConfig().Address), zap.Error(err))
        return
    }

    server := grpc.NewServer()
    cloudFileServer, err := rpc.GetCloudFileServer(config.GetConfig().UserServerAddress)
    if err != nil {
        mlog.Logger().Error("init cloud file server failed", zap.Error(err))
        return
    }
    rpc_impl.RegisterICloudFileServer(server, cloudFileServer)

    mlog.Logger().Info("> Listening at : "+ config.GetConfig().Address)

    // Serve is blocked
    err = server.Serve(listener)
    if err != nil {
        mlog.Logger().Error(fmt.Sprintf("serve on %s failed", config.GetConfig().Address), zap.Error(err))
        return
    }
}
