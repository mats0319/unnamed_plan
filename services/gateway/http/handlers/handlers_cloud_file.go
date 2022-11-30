package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func getCloudFileClientAndConnTarget() (rpc_impl.ICloudFileClient, string, error) {
	conn, err := rce.GetClientConn(mconst.UID_Service_Cloud_File)
	if err != nil {
		mlog.Logger().Error("get client conn failed", zap.Error(err))
		return nil, "", err
	}

	return rpc_impl.NewICloudFileClient(conn), conn.Target(), nil
}

func ListCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.List(context.Background(), &rpc_impl.CloudFile_ListReq{
		Rule:       rpc_impl.CloudFile_ListRule(params.Rule),
		OperatorId: params.OperatorID,
		Page: &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeListCloudFileByUploaderRes(res.Total, filesRPCToHTTP(res.Files...)))
}

func UploadCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.UploadCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Upload(context.Background(), &rpc_impl.CloudFile_UploadReq{
		OperatorId:       params.OperatorID,
		File:             params.File,
		FileName:         params.FileName,
		ExtensionName:    params.ExtensionName,
		FileSize:         params.FileHeader.Size,
		LastModifiedTime: int64(params.LastModifiedTime),
		IsPublic:         params.IsPublic,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Modify(context.Background(), &rpc_impl.CloudFile_ModifyReq{
		OperatorId:       params.OperatorID,
		FileId:           params.FileID,
		Password:         params.Password,
		FileName:         params.FileName,
		ExtensionName:    params.ExtensionName,
		IsPublic:         params.IsPublic,
		File:             params.File,
		FileSize:         params.FileSize,
		LastModifiedTime: int64(params.LastModifiedTime),
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func DeleteCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.DeleteCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Delete(context.Background(), &rpc_impl.CloudFile_DeleteReq{
		OperatorId: params.OperatorID,
		Password:   params.Password,
		FileId:     params.FileID,
	})
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func filesRPCToHTTP(data ...*rpc_impl.CloudFile_Data) []*structure.FileRes {
	res := make([]*structure.FileRes, 0, len(data))
	for i := range data {
		res = append(res, &structure.FileRes{
			FileID:           data[i].FileId,
			FileName:         data[i].FileName,
			LastModifiedTime: time.Duration(data[i].LastModifiedTime),
			FileURL:          data[i].FileUrl,
			IsPublic:         data[i].IsPublic,
			UpdateTime:       time.Duration(data[i].UpdateTime),
			CreatedTime:      time.Duration(data[i].CreatedTime),
		})
	}

	return res
}
