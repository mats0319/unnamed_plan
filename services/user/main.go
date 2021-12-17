package main

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/unnamed_plan/services/user/config"
    "github.com/mats9693/unnamed_plan/services/user/rpc"
    "google.golang.org/grpc"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", config.GetConfig().Address)
    if err != nil {
        fmt.Printf("listen on %s failed, error: %v\n", config.GetConfig().Address, err)
        return
    }

    server := grpc.NewServer()
    rpc_impl.RegisterIUserServer(server, rpc.GetUserServer())

    fmt.Println("> Listening at :", config.GetConfig().Address)

    // Serve is blocked
    err = server.Serve(listener)
    if err != nil {
        fmt.Printf("serve on %s failed, error: %v\n", config.GetConfig().Address, err)
        return
    }
}
