package client

import (
	"context"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"google.golang.org/grpc"
)

type cloudFileClient struct {
	client rpc_impl.ICloudFileClient
}

var _ rpc_impl.ICloudFileClient = (*cloudFileClient)(nil)

var cloudFileClientIns = &cloudFileClient{}

func ConnectCloudFileServer(target string) (*cloudFileClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cloudFileClientIns.client = rpc_impl.NewICloudFileClient(conn)

	return cloudFileClientIns, nil
}

func (c *cloudFileClient) ListByUploader(ctx context.Context, in *rpc_impl.CloudFile_ListByUploaderReq, opts ...grpc.CallOption) (*rpc_impl.CloudFile_ListByUploaderRes, error) {
	return c.client.ListByUploader(ctx, in, opts...)
}

func (c *cloudFileClient) ListPublic(ctx context.Context, in *rpc_impl.CloudFile_ListPublicReq, opts ...grpc.CallOption) (*rpc_impl.CloudFile_ListPublicRes, error) {
	return c.client.ListPublic(ctx, in, opts...)
}

func (c *cloudFileClient) Upload(ctx context.Context, in *rpc_impl.CloudFile_UploadReq, opts ...grpc.CallOption) (*rpc_impl.CloudFile_UploadRes, error) {
	return c.client.Upload(ctx, in, opts...)
}

func (c *cloudFileClient) Modify(ctx context.Context, in *rpc_impl.CloudFile_ModifyReq, opts ...grpc.CallOption) (*rpc_impl.CloudFile_ModifyRes, error) {
	return c.client.Modify(ctx, in, opts...)
}

func (c *cloudFileClient) Delete(ctx context.Context, in *rpc_impl.CloudFile_DeleteReq, opts ...grpc.CallOption) (*rpc_impl.CloudFile_DeleteRes, error) {
	return c.client.Delete(ctx, in, opts...)
}
