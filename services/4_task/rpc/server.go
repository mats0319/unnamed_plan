package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/4_task/config"
	"github.com/mats9693/unnamed_plan/services/4_task/db"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded/invoke"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
)

type taskServerImpl struct {
	rpc_impl.UnimplementedITaskServer
}

var taskServerImplIns = &taskServerImpl{}

func GetTaskServer() rpc_impl.ITaskServer {
	return taskServerImplIns
}

func (t *taskServerImpl) List(_ context.Context, req *rpc_impl.Task_ListReq) (*rpc_impl.Task_ListRes, error) {
	res := &rpc_impl.Task_ListRes{}

	if len(req.OperatorId) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams, zap.String("operator", req.OperatorId))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	tasks, count, err := db.GetTaskDao().QueryByPoster(req.OperatorId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	if len(tasks) > config.GetConfig().MaxRecords {
		tasks = tasks[:config.GetConfig().MaxRecords]
	}

	res.Total = uint32(count)
	res.Tasks = tasksDBToRPC(tasks...)

	return res, nil
}

func (t *taskServerImpl) Create(_ context.Context, req *rpc_impl.Task_CreateReq) (*rpc_impl.Task_CreateRes, error) {
	res := &rpc_impl.Task_CreateRes{}

	if len(req.OperatorId) < 1 || len(req.TaskName) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("task name", req.TaskName))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	err := db.GetTaskDao().Insert(&model.Task{
		TaskName:    req.TaskName,
		PostedBy:    req.OperatorId,
		Description: req.Description,
		PreTaskIDs:  req.PreTaskIds,
		Status:      mconst.TaskStatus_Posted,
	})
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func (t *taskServerImpl) Modify(ctx context.Context, req *rpc_impl.Task_ModifyReq) (*rpc_impl.Task_ModifyRes, error) {
	res := &rpc_impl.Task_ModifyRes{}

	if len(req.OperatorId) < 1 || len(req.TaskId) < 1 || len(req.TaskName) < 1 {
		mlog.Logger().Error(mconst.Error_InvalidParams,
			zap.String("operator", req.OperatorId),
			zap.String("task", req.TaskId),
			zap.String("task name", req.TaskName))
		res.Err = utils.Error_InvalidParams.ToRPC()
		return res, nil
	}

	rpcErr := rce_invoke.AuthUserInfo(ctx, req.OperatorId, req.Password)
	if rpcErr != nil {
		mlog.Logger().Error("auth user info failed", zap.String("error", rpcErr.String()))
		res.Err = rpcErr
		return res, nil
	}

	task, err := db.GetTaskDao().QueryOne(req.TaskId)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	if task.PostedBy != req.OperatorId {
		mlog.Logger().Error(mconst.Error_ModifyOthersTask)
		res.Err = utils.NewExecError(mconst.Error_ModifyOthersTask).ToRPC()
		return res, nil
	}

	if task.TaskName == req.TaskName && task.Description == req.Description &&
		utils.CompareOnStringSliceNotStrict(task.PreTaskIDs, req.PreTaskIds) &&
		uint32(task.Status) == req.Status {
		mlog.Logger().Error(mconst.Error_NoValidModification)
		res.Err = utils.NewExecError(mconst.Error_NoValidModification).ToRPC()
		return res, nil
	}

	updateColumns := make([]string, 0)
	if task.TaskName != req.TaskName {
		task.TaskName = req.TaskName
		updateColumns = append(updateColumns, model.Task_TaskName)
	}
	if task.Description != req.Description {
		task.Description = req.Description
		updateColumns = append(updateColumns, model.Task_Description)
	}
	if !utils.CompareOnStringSliceNotStrict(task.PreTaskIDs, req.PreTaskIds) {
		task.PreTaskIDs = req.PreTaskIds
		updateColumns = append(updateColumns, model.Task_PreTaskIDs)
	}
	if uint32(task.Status) != req.Status {
		task.Status = mconst.TaskStatus(req.Status)
		updateColumns = append(updateColumns, model.Task_Status)
	}

	err = db.GetTaskDao().UpdateColumnsByTaskID(task, updateColumns...)
	if err != nil {
		mlog.Logger().Error(mconst.Error_DBError, zap.Error(err))
		res.Err = utils.NewDBError(err.Error()).ToRPC()
		return res, nil
	}

	return res, nil
}

func tasksDBToRPC(data ...*model.Task) []*rpc_impl.Task_Data {
	res := make([]*rpc_impl.Task_Data, 0, len(data))
	for i := range data {
		res = append(res, &rpc_impl.Task_Data{
			TaskId:      data[i].ID,
			TaskName:    data[i].TaskName,
			Description: data[i].Description,
			PreTaskIds:  data[i].PreTaskIDs,
			Status:      uint32(data[i].Status),
			UpdateTime:  int64(data[i].UpdateTime),
			CreatedTime: int64(data[i].CreatedTime),
		})
	}

	return res
}
