package rpc

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/client"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"os"
)

const uid_RPCClient = "1cd10cb8-ecf5-4855-a886-76b148ed104a"

type rpcClient struct {
	conf *rpcClientConfig

	UserClient      rpc_impl.IUserClient
	CloudFileClient rpc_impl.ICloudFileClient
	NoteClient      rpc_impl.INoteClient
	TaskClient      rpc_impl.ITaskClient
}

type rpcClientConfig struct {
	UserClientTarget      string `json:"userClientTarget"`
	CloudFileClientTarget string `json:"cloudFileClientTarget"`
	NoteClientTarget      string `json:"noteClientTarget"`
	TaskClientTarget      string `json:"taskClientTarget"`
}

var rpcClientIns = &rpcClient{}

func GetRPCClient() *rpcClient {
	return rpcClientIns
}

func init() {
	byteSlice := mconfig.GetConfig(uid_RPCClient)

	rpcClientConfigIns := &rpcClientConfig{}
	err := json.Unmarshal(byteSlice, rpcClientConfigIns)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", uid_RPCClient), zap.Error(err))
		os.Exit(-1)
	}

	rpcClientIns.conf = rpcClientConfigIns

	userClient, err := client.ConnectUserServer(rpcClientIns.conf.UserClientTarget)
	cloudFileClient, err2 := client.ConnectCloudFileServer(rpcClientIns.conf.CloudFileClientTarget)
	noteClient, err3 := client.ConnectNoteServer(rpcClientIns.conf.NoteClientTarget)
	taskClient, err4 := client.ConnectTaskServer(rpcClientIns.conf.TaskClientTarget)

	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		mlog.Logger().Error("establish connection with services failed",
			zap.String("err msg", utils.ErrorsToString(err, err2, err3)))
		os.Exit(-1)
	}

	rpcClientIns.UserClient = userClient
	rpcClientIns.CloudFileClient = cloudFileClient
	rpcClientIns.NoteClient = noteClient
	rpcClientIns.TaskClient = taskClient

	mlog.Logger().Info("> RPC client init finish.")
}
