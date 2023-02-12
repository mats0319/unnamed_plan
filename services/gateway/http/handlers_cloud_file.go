package http

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"net/http"
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
	res := &rpc_impl.CloudFile_ListRes{}

	params := &structure.ListCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.List(context.Background(), (*rpc_impl.CloudFile_ListReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func UploadCloudFile(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.CloudFile_UploadRes{}

	params := &structure.UploadCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Upload(context.Background(), (*rpc_impl.CloudFile_UploadReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func ModifyCloudFile(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.CloudFile_ModifyRes{}

	params := &structure.ModifyCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		res.Err = utils.Error_InvalidParams
		return mhttp.NewResponseDataWithError(res)
	}

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Modify(context.Background(), (*rpc_impl.CloudFile_ModifyReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}

func DeleteCloudFile(r *http.Request) *mhttp.ResponseData {
	res := &rpc_impl.CloudFile_DeleteRes{}

	params := &structure.DeleteCloudFileReqParams{}
	params.Decode(r)

	client, target, err := getCloudFileClientAndConnTarget()
	if err != nil {
		mlog.Logger().Error("get cloud file client failed", zap.Error(err))
		res.Err = utils.NewGetClientError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}

	res, err = client.Delete(context.Background(), (*rpc_impl.CloudFile_DeleteReq)(params))
	if err != nil {
		rce.ReportInvalidTarget(mconst.UID_Service_Cloud_File, target)
		mlog.Logger().Error(mconst.Error_GrpcConnectionError, zap.Error(err))
		res.Err = utils.NewGrpcConnectionError(err.Error())
		return mhttp.NewResponseDataWithError(res)
	}
	if res != nil && res.Err != nil {
		mlog.Logger().Error(mconst.Error_ExecutionError, zap.String("error", res.Err.String()))
		return mhttp.NewResponseDataWithError(res)
	}

	return mhttp.NewResponseData(res)
}
