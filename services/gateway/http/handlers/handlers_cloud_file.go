package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	"github.com/mats9693/unnamed_plan/shared/proto/impl"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
	"time"
)

func ListCloudFileByUploader(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListCloudFileByUploaderReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().CloudFileClient.ListByUploader(context.Background(), &rpc_impl.CloudFile_ListByUploaderReq{
		OperatorId: params.OperatorID,
		Page:       &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListCloudFileByUploaderRes(int(res.Total), filesRPCToHTTP(res.Files...)))
}

func ListPublicCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListPublicCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().CloudFileClient.ListPublic(context.Background(), &rpc_impl.CloudFile_ListPublicReq{
		OperatorId: params.OperatorID,
		Page:       &rpc_impl.Pagination{
			PageSize: uint32(params.PageSize),
			PageNum:  uint32(params.PageNum),
		},
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListPublicCloudFileRes(int(res.Total), filesRPCToHTTP(res.Files...)))
}

func UploadCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.UploadCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().CloudFileClient.Upload(context.Background(), &rpc_impl.CloudFile_UploadReq{
		OperatorId:       params.OperatorID,
		File:             params.File,
		FileName:         params.FileName,
		ExtensionName:    params.ExtensionName,
		FileSize:         params.FileHeader.Size,
		LastModifiedTime: int64(params.LastModifiedTime),
		IsPublic:         params.IsPublic,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response("")
}

func ModifyCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().CloudFileClient.Modify(context.Background(), &rpc_impl.CloudFile_ModifyReq{
		OperatorId:       params.OperatorID,
		FileId:           params.FileID,
		Password:         params.Password,
		FileName:         params.FileName,
		ExtensionName:    params.ExtensionName,
		IsPublic:         params.IsPublic,
		File:             params.File,
		FileSize:         params.FileHeader.Size,
		LastModifiedTime: int64(params.LastModifiedTime),
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response("")
}

func DeleteCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.DeleteCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().CloudFileClient.Delete(context.Background(), &rpc_impl.CloudFile_DeleteReq{
		OperatorId: params.OperatorID,
		Password:   params.Password,
		FileId:     params.FileID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeDeleteCloudFileRes(true))
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
