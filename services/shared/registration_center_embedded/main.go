package rce

import (
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
)

type rce struct {
	instance *rceImpl

	targetMap sync.Map // service id - *rpcTarget
}

var rceIns = &rce{
	instance: &rceImpl{},
}

// Init for not-business service but need to invoke business services, like 'gateway'
func Init(registrationCenterTarget string) error {
	err := rceIns.instance.initialize(registrationCenterTarget)
	if err != nil {
		mlog.Logger().Error("init rc embedded failed", zap.Error(err))
		return err
	}

	return nil
}

// InitAndRegister for business service
func InitAndRegister(registrationCenterTarget string, serviceID string, target string) (rpc_impl.IRegistrationCenterEmbeddedServer, error) {
	err := Init(registrationCenterTarget)
	if err != nil {
		mlog.Logger().Error("init rc embedded failed", zap.Error(err))
		return nil, err
	}

	err = rceIns.instance.register(serviceID, target)
	if err != nil {
		mlog.Logger().Error("register service failed", zap.Error(err))
		return nil, err
	}

	return rceIns.instance, nil
}

func GetClientConn(serviceID string) (*grpc.ClientConn, error) {
	if !rceIns.instance.init {
		return nil, errors.New("RCE module not init")
	}

	rpcTargetI, ok := rceIns.targetMap.Load(serviceID)
	rpcTargetIns, _ := rpcTargetI.(*rpcTarget) // make sure value type of map is '*rpcTarget'
	if !ok || len(rpcTargetIns.list) < 1 {     // key is not exist or value is empty
		targetList, err := rceIns.instance.ListServiceTarget(serviceID)
		if err != nil {
			mlog.Logger().Error("list service target failed", zap.Error(err))
			return nil, err
		}

		rpcTargetIns = newRPCTarget(targetList)

		rceIns.targetMap.Store(serviceID, rpcTargetIns)
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

// ReportInvalidTarget if there are something wrong when service 'A' invoke service 'B',
// and error is 'connection error' ('err' field support by grpc is not nil).
// We consider that 'target' is invalid, and service 'A' report it to Registration Center Embedded module,
// this is response function of RCE module.
//
// in this version, RCE module delete invalid target directly
func ReportInvalidTarget(serviceID string, target string) {
	rpcTargetI, ok := rceIns.targetMap.Load(serviceID)
	rpcTargetIns, _ := (rpcTargetI).(*rpcTarget)
	if !ok || len(rpcTargetIns.list) < 1 { // unexpected data type or empty data
		return
	}

	index := -1
	for i := 0; i < len(rpcTargetIns.list); i++ {
		if rpcTargetIns.list[i] == target {
			index = i
			break
		}
	}

	if index < 0 { // target not exist
		return
	}

	lastIndex := len(rpcTargetIns.list) - 1

	// del 'target' and reset polling index to 0
	rpcTargetIns.list[index], rpcTargetIns.list[lastIndex] = rpcTargetIns.list[lastIndex], rpcTargetIns.list[index]
	rpcTargetIns.list = rpcTargetIns.list[:lastIndex]

	rpcTargetIns.index = 0

	rceIns.targetMap.Store(serviceID, rpcTargetIns)
}
