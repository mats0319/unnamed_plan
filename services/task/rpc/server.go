package rpc

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/mats9693/unnamed_plan/services/task/config"
	"github.com/mats9693/unnamed_plan/services/task/db"
)

type taskServerImpl struct {
	rpc_impl.UnimplementedITaskServer
}

var _ rpc_impl.ITaskServer = (*taskServerImpl)(nil)

var taskServerImplIns = &taskServerImpl{}

func GetTaskServer() *taskServerImpl {
	return taskServerImplIns
}

func (t *taskServerImpl) List(_ context.Context, req *rpc_impl.Task_ListReq) (*rpc_impl.Task_ListRes, error) {
	if len(req.OperatorId) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	tasks, count, err := db.GetTaskDao().QueryByPoster(req.OperatorId)
	if err != nil {
		return nil, err
	}

	if len(tasks) > config.GetConfig().MaxRecords {
		tasks = tasks[:config.GetConfig().MaxRecords]
	}

	return &rpc_impl.Task_ListRes{
		Total: uint32(count),
		Tasks: tasksDBToRPC(tasks...),
	}, nil
}

func (t *taskServerImpl) Create(_ context.Context, req *rpc_impl.Task_CreateReq) (*rpc_impl.Task_CreateRes, error) {
	if len(req.OperatorId) < 1 || len(req.TaskName) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	err := db.GetTaskDao().Insert(&model.Task{
		TaskName:    req.TaskName,
		PostedBy:    req.OperatorId,
		Description: req.Description,
		PreTaskIDs:  req.PreTaskIds,
		Status:      mconst.TaskStatus_Posted,
	})
	if err != nil {
		return nil, err
	}

	return &rpc_impl.Task_CreateRes{}, nil
}

func (t *taskServerImpl) Modify(_ context.Context, req *rpc_impl.Task_ModifyReq) (*rpc_impl.Task_ModifyRes, error) {
	if len(req.OperatorId) < 1 || len(req.TaskId) < 1 || len(req.TaskName) < 1 {
		return nil, utils.NewError(mconst.Error_InvalidParams)
	}

	task, err := db.GetTaskDao().QueryOne(req.TaskId)
	if err != nil {
		return nil, err
	}

	if task.PostedBy != req.OperatorId {
		return nil, utils.NewError(mconst.Error_ModifyOthersTask)
	}

	if task.TaskName == req.TaskName && task.Description == req.Description &&
		utils.CompareOnStringSliceNotStrict(task.PreTaskIDs, req.PreTaskIds) && uint32(task.Status) == req.Status {
		return nil, utils.NewError(mconst.Error_NoValidModification)
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
		return nil, err
	}

	return &rpc_impl.Task_ModifyRes{}, nil
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
