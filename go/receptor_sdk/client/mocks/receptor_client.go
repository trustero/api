package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Receptor struct {
	mock.Mock
}

func (t *Receptor) GetConfiguration(ctx context.Context, request *receptor_v1.ReceptorOID, opts ...grpc.CallOption) (*receptor_v1.ReceptorConfiguration, error) {
	args := t.Called(ctx, request)
	return args.Get(0).(*receptor_v1.ReceptorConfiguration), args.Error(1)
}

func (t *Receptor) Verified(ctx context.Context, request *receptor_v1.Credential, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := t.Called(ctx, request)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (t *Receptor) Discovered(ctx context.Context, request *receptor_v1.Services, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	args := t.Called(ctx, request)
	return args.Get(0).(*wrapperspb.StringValue), args.Error(1)
}
func (t *Receptor) Report(ctx context.Context, request *receptor_v1.Finding, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	args := t.Called(ctx, request)
	return args.Get(0).(*wrapperspb.StringValue), args.Error(1)
}

func (t *Receptor) Notify(ctx context.Context, request *receptor_v1.JobResult, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := t.Called(ctx, request)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}
