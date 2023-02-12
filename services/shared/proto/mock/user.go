package mock_rpc_impl

import (
	"bou.ke/monkey"
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded/invoke"
)

func MockUserAuthenticate() {
	monkey.Patch(rce_invoke.AuthUserInfo, func(_ context.Context, _ *rpc_impl.User_AuthenticateReq) *rpc_impl.Error {
		return nil
	})
}
