package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/2_cloud_file/config"
	"github.com/mats9693/unnamed_plan/services/2_cloud_file/db"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded/invoke"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

const backupFileSuffix = ".backup"

type cloudFileServerImpl struct {
	rpc_impl.UnimplementedICloudFileServer
}

var cloudFileServerImplIns = &cloudFileServerImpl{}

func GetCloudFileServer() rpc_impl.ICloudFileServer {
	return cloudFileServerImplIns
}

func (s *cloudFileServerImpl) List(_ context.Context, req *rpc_impl.CloudFile_ListReq) (*rpc_impl.CloudFile_ListRes, error) {
	res := &rpc_impl.CloudFile_ListRes{}

	if len(req.OperatorId) < 1 || req.Page == nil || req.Page.PageSize < 1 || req.Page.PageNum < 1 ||
		(req.Rule != rpc_impl.CloudFile_UPLOADER && req.Rule != rpc_impl.CloudFile_PUBLIC) {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("page info", req.Page.String()),
			zap.Int32("query condition", int32(req.Rule)))
		res.Err = utils.Error_InvalidParams
		return res, nil
	}

	pageSize := int(req.Page.PageSize)
	pageNum := int(req.Page.PageNum)

	var (
		files []*model.CloudFile
		count int
		err   error
	)
	switch req.Rule {
	case rpc_impl.CloudFile_UPLOADER:
		files, count, err = db.GetCloudFileDao().QueryPageByUploader(pageSize, pageNum, req.OperatorId)
	case rpc_impl.CloudFile_PUBLIC:
		files, count, err = db.GetCloudFileDao().QueryPageInPublic(pageSize, pageNum, req.OperatorId)
	}
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error())
		return res, nil
	}

	res.Total = uint32(count)
	res.Files = filesDBToRPC(files...)

	return res, nil
}

func (s *cloudFileServerImpl) Upload(ctx context.Context, req *rpc_impl.CloudFile_UploadReq) (*rpc_impl.CloudFile_UploadRes, error) {
	res := &rpc_impl.CloudFile_UploadRes{}

	if len(req.OperatorId) < 1 || len(req.FileName) < 1 || len(req.ExtensionName) < 1 ||
		req.FileSize < 1 || req.LastModifiedTime < 1 || len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("file name", req.FileName),
			zap.String("extension name", req.ExtensionName),
			zap.Int64("file size", req.FileSize),
			zap.Int64("last modified time", req.LastModifiedTime),
			zap.String("password", req.Password))
		res.Err = utils.Error_InvalidParams
		return res, nil
	}

	rpcErr := rce_invoke.AuthUserInfo(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	dir := spliceFileDir(req.IsPublic, req.OperatorId)
	// make sure private path exist, public path has checked in init step
	if !req.IsPublic {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
			res.Err = utils.NewExecError(err.Error())
			return res, nil
		}
	}

	// save file
	fileID := utils.CalcSHA256(req.OperatorId + time.Now().GoString())
	absolutePath := spliceFilePath(dir, fileID, req.ExtensionName)

	err := os.WriteFile(absolutePath, req.File, 0755)
	if err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
		res.Err = utils.NewExecError(err.Error())
		return res, nil
	}

	// save db, if failed, remove file
	err = db.GetCloudFileDao().Insert(&model.CloudFile{
		UploadedBy:       req.OperatorId,
		FileID:           fileID,
		FileName:         req.FileName,
		ExtensionName:    req.ExtensionName,
		LastModifiedTime: time.Duration(req.LastModifiedTime),
		FileSize:         req.FileSize,
		IsPublic:         req.IsPublic,
	})
	if err != nil {
		err2 := os.Remove(absolutePath)
		mlog.Logger().Error(mconst.Error_ExecutionError,
			zap.NamedError(mconst.Error_DBError, err),
			zap.NamedError(mconst.Error_ExecutionError, err2))
		res.Err = utils.NewExecError(utils.ErrorsToString(err, err2))
		return res, nil
	}

	return res, nil
}

func (s *cloudFileServerImpl) Modify(ctx context.Context, req *rpc_impl.CloudFile_ModifyReq) (*rpc_impl.CloudFile_ModifyRes, error) {
	res := &rpc_impl.CloudFile_ModifyRes{}

	if len(req.OperatorId) < 1 || len(req.FileId) < 1 || len(req.Password) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("file", req.FileId),
			zap.String("password", req.Password))
		res.Err = utils.Error_InvalidParams
		return res, nil
	}

	rpcErr := rce_invoke.AuthUserInfo(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	fileRecord, err := db.GetCloudFileDao().QueryOne(req.FileId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error())
		return res, nil
	}

	if len(req.FileName)+len(req.ExtensionName) < 1 && fileRecord.IsPublic == req.IsPublic && len(req.File) < 1 {
		mlog.Logger().Error(mconst.Error_NoValidModification)
		res.Err = utils.Error_NoValidModification
		return res, nil
	}

	updateColumns := make([]string, 0, 5)

	// update file, backup old file in private folder
	if len(req.File) > 0 {
		// make sure private dir exist, only when old file is public
		if fileRecord.IsPublic {
			if err = os.MkdirAll(spliceFileDir(false, req.OperatorId), 0755); err != nil {
				mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
				res.Err = utils.NewExecError(err.Error())
				return res, nil
			}
		}

		// backup old file, may implicitly move file from public path to private path
		oldPath := spliceFilePath(spliceFileDir(fileRecord.IsPublic, req.OperatorId), fileRecord.FileID, req.ExtensionName)
		privPath := spliceFilePath(spliceFileDir(false, req.OperatorId), fileRecord.FileID, req.ExtensionName)

		err = os.Rename(oldPath, getValidBackupFileName(privPath))
		if err != nil {
			mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
			res.Err = utils.NewExecError(err.Error())
			return res, nil
		}

		// save new file
		absolutePath := spliceFilePath(spliceFileDir(req.IsPublic, req.OperatorId), fileRecord.FileID, req.ExtensionName)
		err = os.WriteFile(absolutePath, req.File, 0755)
		if err != nil {
			mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))
			res.Err = utils.NewExecError(err.Error())
			return res, nil
		}

		// update db record: file size and file last modified time
		fileRecord.FileSize = req.FileSize
		updateColumns = append(updateColumns, model.CloudFile_FileSize)

		fileRecord.LastModifiedTime = time.Duration(req.LastModifiedTime)
		updateColumns = append(updateColumns, model.CloudFile_LastModifiedTime)
	}

	if len(req.FileName) > 0 {
		fileRecord.FileName = req.FileName
		updateColumns = append(updateColumns, model.CloudFile_FileName)
	}
	if len(req.ExtensionName) > 0 && fileRecord.ExtensionName != req.ExtensionName {
		fileRecord.ExtensionName = req.ExtensionName
		updateColumns = append(updateColumns, model.CloudFile_ExtensionName)
	}
	if fileRecord.IsPublic != req.IsPublic {
		fileRecord.IsPublic = req.IsPublic
		updateColumns = append(updateColumns, model.CloudFile_IsPublic)
	}

	// update db record, if update failed and have new file, remove file
	err = db.GetCloudFileDao().UpdateColumnsByFileID(fileRecord, updateColumns...)
	if err != nil {
		errMsg := err.Error()
		if len(req.File) > 0 {
			dir := spliceFileDir(req.IsPublic, req.OperatorId)
			absolutePath := spliceFilePath(dir, fileRecord.FileID, req.ExtensionName)
			err2 := os.Remove(absolutePath)
			if err2 != nil {
				errMsg += err2.Error()
				mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err2))
			}
		}

		mlog.Logger().Error(mconst.Error_ExecutionError, zap.Error(err))

		res.Err = utils.NewExecError(errMsg)
		return res, nil
	}

	return res, nil
}

func (s *cloudFileServerImpl) Delete(ctx context.Context, req *rpc_impl.CloudFile_DeleteReq) (*rpc_impl.CloudFile_DeleteRes, error) {
	res := &rpc_impl.CloudFile_DeleteRes{}

	if len(req.OperatorId) < 1 || len(req.Password) < 1 || len(req.FileId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("password", req.Password),
			zap.String("file", req.FileId))
		res.Err = utils.Error_InvalidParams
		return res, nil
	}

	rpcErr := rce_invoke.AuthUserInfo(ctx, &rpc_impl.User_AuthenticateReq{
		UserId:   req.OperatorId,
		Password: req.Password,
	})
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	fileRecord, err := db.GetCloudFileDao().QueryOne(req.FileId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error())
		return res, nil
	}

	fileRecord.IsDeleted = true

	err = db.GetCloudFileDao().UpdateColumnsByFileID(fileRecord, model.CloudFile_IsDeleted)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error())
		return res, nil
	}

	return res, nil
}

func getValidBackupFileName(absolutePath string) string {
	absolutePath += backupFileSuffix

	for i := 0; ; i++ {
		fileInfo, err := os.Stat(absolutePath + strconv.Itoa(i))
		if !(err == nil && !fileInfo.IsDir()) {
			absolutePath += strconv.Itoa(i)
			break
		}
	}

	return absolutePath
}

func filesDBToRPC(data ...*model.CloudFile) []*rpc_impl.CloudFile_Data {
	res := make([]*rpc_impl.CloudFile_Data, 0, len(data))
	for i := range data {
		res = append(res, &rpc_impl.CloudFile_Data{
			FileId:           data[i].FileID,
			FileName:         data[i].FileName,
			LastModifiedTime: int64(data[i].LastModifiedTime),
			FileUrl:          spliceFilePath(spliceFileURL(data[i].IsPublic, data[i].UploadedBy), data[i].FileID, data[i].ExtensionName),
			IsPublic:         data[i].IsPublic,
			UpdateTime:       int64(data[i].UpdateTime),
			CreatedTime:      int64(data[i].CreatedTime),
		})
	}

	return res
}

func spliceFileURL(isPublic bool, operatorID string) string {
	url := ""

	if isPublic {
		url += config.GetConfig().CloudFilePublicDir
	} else {
		url += operatorID
	}

	return utils.FormatDirSuffix(url)
}

func spliceFileDir(isPublic bool, operatorID string) string {
	dir := utils.FormatDirSuffix(config.GetConfig().CloudFileRootPath)
	dir += spliceFileURL(isPublic, operatorID)

	return utils.FormatDirSuffix(dir)
}

// spliceFilePath return absolute path
func spliceFilePath(dir string, fileName string, extensionName string) string {
	return dir + fileName + "." + extensionName
}
