package main

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/cloud_file/config"
	"github.com/mats9693/unnamed_plan/services/cloud_file/db"
	"github.com/mats9693/unnamed_plan/services/cloud_file/rpc"
	mdb "github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
	"strings"
)

func main() {
	initialize.Init("config.json", mdb.Init, config.Init, db.Init, InitCloudFileDir)

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

	mlog.Logger().Info("> Listening at : " + config.GetConfig().Address)

	// Serve is blocked
	err = server.Serve(listener)
	if err != nil {
		mlog.Logger().Error(fmt.Sprintf("serve on %s failed", config.GetConfig().Address), zap.Error(err))
		return
	}
}

func InitCloudFileDir() {
	root := config.GetConfig().CloudFileRootPath
	if len(root) < 1 {
		executableAbsolutePath, err := os.Executable()
		if err != nil {
			fmt.Println("get executable failed, error:", err)
			os.Exit(-1)
		}
		executableAbsolutePath = strings.ReplaceAll(executableAbsolutePath, "\\", "/")
		executableDir := utils.FormatDirSuffix(path.Dir(executableAbsolutePath))

		root = executableDir + "files/"
	}
	cloudFileDir := utils.FormatDirSuffix(root) + config.GetConfig().CloudFilePublicDir

	err := os.MkdirAll(cloudFileDir, 0755)
	if err != nil {
		mlog.Logger().Error("os.MkdirAll failed", zap.Error(err))
		os.Exit(-1)
	}

	mlog.Logger().Info("> Cloud file directory init finish.")
}
