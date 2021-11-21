package main

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/cloud_file/config"
    "github.com/mats9693/unnamed_plan/services/cloud_file/rpc"
    "github.com/mats9693/unnamed_plan/shared/proto/impl"
    "github.com/mats9693/utils/toy_server/utils"
    "google.golang.org/grpc"
    "net"
    "os"
)

func init() {
    root := config.GetConfig().CloudFileRootPath
    path := mutils.FormatDirSuffix(root) + config.GetConfig().CloudFilePublicDir

    err := os.MkdirAll(path, 0755)
    if err != nil {
        fmt.Println("os.MkdirAll failed, error:", err.Error())
        os.Exit(-1)
    }

    fmt.Println("> Cloud file directory init finish.")
}

func main() {
    listener, err := net.Listen("tcp", config.GetConfig().Address)
    if err != nil {
        fmt.Printf("listen on %s failed, error: %v\n", config.GetConfig().Address, err)
        return
    }

    server := grpc.NewServer()
    cloudFileServer, err := rpc.GetCloudFileServer(config.GetConfig().UserServerAddress)
    if err != nil {
        fmt.Println("init thinking note server failed, error:", err)
        return
    }
    rpc_impl.RegisterICloudFileServer(server, cloudFileServer)

    fmt.Println("> Listening at :", config.GetConfig().Address)

    // Serve is blocked
    err = server.Serve(listener)
    if err != nil {
        fmt.Printf("serve on %s failed, error: %v\n", config.GetConfig().Address, err)
        return
    }
}
