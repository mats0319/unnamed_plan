package handlers

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/gateway/http/structure_defination"
	"github.com/mats9693/unnamed_plan/services/gateway/rpc"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func ListTask(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	res, err := rpc.GetRPCClient().TaskClient.List(context.Background(), &rpc_impl.Task_ListReq{
		OperatorId: params.OperatorID,
	})
	if err != nil {
		mlog.Logger().Error("list task failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeListTaskRes(res.Total, tasksRPCToHTTP(res.Tasks...)))
}

func CreateTask(r *http.Request) *mhttp.ResponseData {
	params := &structure.CreateTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().TaskClient.Create(context.Background(), &rpc_impl.Task_CreateReq{
		OperatorId:  params.OperatorID,
		TaskName:    params.TaskName,
		Description: params.Description,
		PreTaskIds:  params.PreTaskIDs,
	})
	if err != nil {
		mlog.Logger().Error("create task failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(mconst.EmptyHTTPRes)
}

func ModifyTask(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyTaskReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		mlog.Logger().Error("parse request params failed", zap.String("err msg", errMsg))
		return mhttp.ResponseWithError(errMsg)
	}

	_, err := rpc.GetRPCClient().TaskClient.Modify(context.Background(), &rpc_impl.Task_ModifyReq{
		OperatorId:  params.OperatorID,
		TaskId:      params.TaskID,
		Password:    params.Password,
		TaskName:    params.TaskName,
		Description: params.Description,
		PreTaskIds:  params.PreTaskIDs,
		Status:      uint32(params.Status),
	})
	if err != nil {
		mlog.Logger().Error("create task failed", zap.Error(err))
		return mhttp.ResponseWithError(err.Error())
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
