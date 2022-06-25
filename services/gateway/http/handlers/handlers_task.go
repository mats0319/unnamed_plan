package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func getTaskClient() (rpc_impl.ITaskClient, error) {
	conn, err := rc_embedded.GetClientConn(mconst.UID_Service_Task)
	if err != nil {
		mlog.Logger().Error("get client conn failed", zap.Error(err))
		return nil, err
	}

	return rpc_impl.NewITaskClient(conn), nil
}

func ListTask(r *http.Request) *mresponse.ResponseData {
	params := &structure.ListTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, err := getTaskClient()
	if err != nil {
		mlog.Logger().Error("get task client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.List(context.Background(), &rpc_impl.Task_ListReq{
		OperatorId: params.OperatorID,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(structure.MakeListTaskRes(res.Total, tasksRPCToHTTP(res.Tasks...)))
}

func CreateTask(r *http.Request) *mresponse.ResponseData {
	params := &structure.CreateTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, err := getTaskClient()
	if err != nil {
		mlog.Logger().Error("get task client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Create(context.Background(), &rpc_impl.Task_CreateReq{
		OperatorId:  params.OperatorID,
		TaskName:    params.TaskName,
		Description: params.Description,
		PreTaskIds:  params.PreTaskIDs,
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyTask(r *http.Request) *mresponse.ResponseData {
	params := &structure.ModifyTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	client, err := getTaskClient()
	if err != nil {
		mlog.Logger().Error("get task client failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	res, err := client.Modify(context.Background(), &rpc_impl.Task_ModifyReq{
		OperatorId:  params.OperatorID,
		TaskId:      params.TaskID,
		Password:    params.Password,
		TaskName:    params.TaskName,
		Description: params.Description,
		PreTaskIds:  params.PreTaskIDs,
		Status:      uint32(params.Status),
	})
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}
	if res != nil && res.Err != nil {
		return mhttp.ResponseWithError(res.Err.String())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func tasksRPCToHTTP(data ...*rpc_impl.Task_Data) []*structure.TaskRes {
	res := make([]*structure.TaskRes, 0, len(data))
	for i := range data {
		res = append(res, &structure.TaskRes{
			TaskID:      data[i].TaskId,
			TaskName:    data[i].TaskName,
			Description: data[i].Description,
			PreTaskIDs:  data[i].PreTaskIds,
			Status:      mconst.TaskStatus(data[i].Status),
			UpdateTime:  time.Duration(data[i].UpdateTime),
			CreatedTime: time.Duration(data[i].CreatedTime),
		})
	}

	return res
}
