package main

import (
    "github.com/mats9693/unnamed_plan/services/config_center/http"
    "github.com/mats9693/unnamed_plan/services/shared/http"
)

func main() {
    mhttp.StartServer(http.GetHandler()) // todo: add 'go' when support rpc service

    //listener, err := net.Listen("tcp", config.GetConfig().Address)
    //if err != nil {
    //    mlog.Logger().Error(fmt.Sprintf("listen on %s failed", config.GetConfig().Address), zap.Error(err))
    //    return
    //}
    //
    //server := grpc.NewServer()
    //rpc_impl.RegisterIUserServer(server, rpc.GetUserServer())
    //
    //mlog.Logger().Info("> Listening at : " + config.GetConfig().Address)
    //
    //// Serve is blocked
    //err = server.Serve(listener)
    //if err != nil {
    //    mlog.Logger().Error(fmt.Sprintf("serve on %s failed", config.GetConfig().Address), zap.Error(err))
    //    return
    //}
}
