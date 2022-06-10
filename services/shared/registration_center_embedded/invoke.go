package rc_embedded

import (
	mconst "github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
)

// rcEmbedded
// Init with rc core target at first,
// Register service instance to rc core, need Init first
// GetClient return service client, need Init first
type rcEmbedded struct {
	registrationCenterImpl *rcEmbeddedImpl

	targetMap sync.Map // service id - *rpcTarget
}

var rcEmbeddedIns = &rcEmbedded{}

func Init(registrationCenterTarget string) {
	conn, err := grpc.Dial(registrationCenterTarget, grpc.WithInsecure())
	if err != nil {
		mlog.Logger().Error("grpc dial failed", zap.Error(err))
		return
	}

	rcEmbeddedIns.registrationCenterImpl = &rcEmbeddedImpl{
		rcCoreClient: rpc_impl.NewIRegistrationCenterCoreClient(conn),
	}
}

func Register(serviceID string, target string) error {
	if rcEmbeddedIns.registrationCenterImpl == nil {
		return errors.New("RC embedded not init")
	}

	err := rcEmbeddedIns.registrationCenterImpl.register(serviceID, target)
	if err != nil {
		mlog.Logger().Error("register service failed", zap.Error(err))
		return err
	}

	return nil
}

func GetUserClient() (rpc_impl.IUserClient, error) {
	conn, err := getClientConn(mconst.UID_Service_User)
	if err != nil {
		mlog.Logger().Error("get grpc client conn failed", zap.Error(err))
		return nil, err
	}

	return rpc_impl.NewIUserClient(conn), nil
}

func GetCloudFileClient() (rpc_impl.ICloudFileClient, error) {
	conn, err := getClientConn(mconst.UID_Service_Cloud_File)
	if err != nil {
		mlog.Logger().Error("get grpc client conn failed", zap.Error(err))
		return nil, err
	}

	return rpc_impl.NewICloudFileClient(conn), nil
}

func GetNoteClient() (rpc_impl.INoteClient, error) {
	conn, err := getClientConn(mconst.UID_Service_Note)
	if err != nil {
		mlog.Logger().Error("get grpc client conn failed", zap.Error(err))
		return nil, err
	}

	return rpc_impl.NewINoteClient(conn), nil
}

func GetTaskClient() (rpc_impl.ITaskClient, error) {
	conn, err := getClientConn(mconst.UID_Service_Task)
	if err != nil {
		mlog.Logger().Error("get grpc client conn failed", zap.Error(err))
		return nil, err
	}

	return rpc_impl.NewITaskClient(conn), nil
}

func getClientConn(serviceID string) (*grpc.ClientConn, error) {
	if rcEmbeddedIns.registrationCenterImpl == nil {
		return nil, errors.New("RC embedded not init")
	}

	targetI, ok := rcEmbeddedIns.targetMap.Load(serviceID)

	rpcTargetIns, _ := targetI.(*rpcTarget) // make sure value type of map is '*rpcTarget'
	if !ok || len(rpcTargetIns.list) < 1 {  // key is not exist or value is empty
		targetList, err := rcEmbeddedIns.registrationCenterImpl.ListServiceTarget(serviceID)
		if err != nil {
			mlog.Logger().Error("list service target failed", zap.Error(err))
			return nil, err
		}

		rpcTargetIns = newTarget(targetList)

		rcEmbeddedIns.targetMap.Store(serviceID, rpcTargetIns)
	}

	target, err := rpcTargetIns.getTarget()
	if err != nil {
		mlog.Logger().Error("get target failed", zap.Error(err))
		return nil, err
	}

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		mlog.Logger().Error("grpc dial failed", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
