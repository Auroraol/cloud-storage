// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: log_service.proto

package sshservicerpc

import (
	"context"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	DeleteSshInfoReq  = pb.DeleteSshInfoReq
	GetSshInfosReq    = pb.GetSshInfosReq
	OperationLogReq   = pb.OperationLogReq
	OperationLogResp  = pb.OperationLogResp
	SshInfoDetailResp = pb.SshInfoDetailResp
	SshInfoListResp   = pb.SshInfoListResp
	SshInfoReq        = pb.SshInfoReq
	SshInfoResp       = pb.SshInfoResp

	SshServiceRpc interface {
		// 保存ssh信息
		SaveSshInfo(ctx context.Context, in *SshInfoReq, opts ...grpc.CallOption) (*SshInfoResp, error)
		// 删除ssh信息
		DeleteSshInfo(ctx context.Context, in *DeleteSshInfoReq, opts ...grpc.CallOption) (*SshInfoResp, error)
		// 查询ssh信息
		GetSshInfo(ctx context.Context, in *GetSshInfosReq, opts ...grpc.CallOption) (*SshInfoListResp, error)
	}

	defaultSshServiceRpc struct {
		cli zrpc.Client
	}
)

func NewSshServiceRpc(cli zrpc.Client) SshServiceRpc {
	return &defaultSshServiceRpc{
		cli: cli,
	}
}

// 保存ssh信息
func (m *defaultSshServiceRpc) SaveSshInfo(ctx context.Context, in *SshInfoReq, opts ...grpc.CallOption) (*SshInfoResp, error) {
	client := pb.NewSshServiceRpcClient(m.cli.Conn())
	return client.SaveSshInfo(ctx, in, opts...)
}

// 删除ssh信息
func (m *defaultSshServiceRpc) DeleteSshInfo(ctx context.Context, in *DeleteSshInfoReq, opts ...grpc.CallOption) (*SshInfoResp, error) {
	client := pb.NewSshServiceRpcClient(m.cli.Conn())
	return client.DeleteSshInfo(ctx, in, opts...)
}

// 查询ssh信息
func (m *defaultSshServiceRpc) GetSshInfo(ctx context.Context, in *GetSshInfosReq, opts ...grpc.CallOption) (*SshInfoListResp, error) {
	client := pb.NewSshServiceRpcClient(m.cli.Conn())
	return client.GetSshInfo(ctx, in, opts...)
}
