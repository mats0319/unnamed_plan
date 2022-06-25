package mock_rpc_impl

import (
	"bou.ke/monkey"
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/client"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
)

func MockUserAuthenticate() {
	monkey.Patch(client.AuthUserInfo, func(ctx context.Context, userID string, password string) *rpc_impl.Error {
		return nil
	})
}
