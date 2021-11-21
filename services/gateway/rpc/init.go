package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/utils"
	"github.com/mats9693/unnamed_plan/shared/proto/client"
	"github.com/mats9693/unnamed_plan/shared/proto/impl"
	"github.com/mats9693/utils/toy_server/config"
	"os"
)

const uid_RPCClient = "1cd10cb8-ecf5-4855-a886-76b148ed104a"

var rpcClientIns = &rpcClient{}

type rpcClient struct {
	conf *rpcClientConfig

	UserClient         rpc_impl.IUserClient
	CloudFileClient    rpc_impl.ICloudFileClient
	ThinkingNoteClient rpc_impl.IThinkingNoteClient
}

type rpcClientConfig struct {
	UserClientTarget         string `json:"userClientTarget"`
	CloudFileClientTarget    string `json:"cloudFileClientTarget"`
	ThinkingNoteClientTarget string `json:"thinkingNoteClientTarget"`
}

func init() {
	byteSlice := mconfig.GetConfig(uid_RPCClient)

	rpcClientConfigIns := &rpcClientConfig{}
	err := json.Unmarshal(byteSlice, rpcClientConfigIns)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", uid_RPCClient, err)
		os.Exit(-1)
	}

	rpcClientIns.conf = rpcClientConfigIns

	userClient, err := client.ConnectUserServer(rpcClientIns.conf.UserClientTarget)
	cloudFileClient, err2 := client.ConnectCloudFileServer(rpcClientIns.conf.CloudFileClientTarget)
	thinkingNoteClient, err3 := client.ConnectThinkingNoteServer(rpcClientIns.conf.ThinkingNoteClientTarget)

	if err != nil || err2 != nil || err3 != nil {
		fmt.Println("establish connection with services failed, error:", utils.ErrorsToString(err, err2, err3))
		os.Exit(-1)
	}

	rpcClientIns.UserClient = userClient
	rpcClientIns.CloudFileClient = cloudFileClient
	rpcClientIns.ThinkingNoteClient = thinkingNoteClient

	fmt.Println("> RPC client init finish.")
}

func GetRPCClient() *rpcClient {
	return rpcClientIns
}
