package mock_rpc_impl

import (
	"bou.ke/monkey"
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/registration_center_embedded/invoke"
)

func MockUserAuthenticate() {
	monkey.Patch(rce_invoke.AuthUserInfo, func(ctx context.Context, userID string, password string) *rpc_impl.Error {
		return nil
	})
}
