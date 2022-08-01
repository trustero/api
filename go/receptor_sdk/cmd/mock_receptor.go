package cmd

import (
	"context"

	receptor "github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type MockReceptorClient struct{}

func (rc *MockReceptorClient) Verified(ctx context.Context, in *receptor.Credential, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (rc *MockReceptorClient) GetConfiguration(ctx context.Context, in *receptor.ReceptorOID, opts ...grpc.CallOption) (*receptor.ReceptorConfiguration, error) {
	return &receptor.ReceptorConfiguration{
		ReceptorObjectId:       "",
		Credential:             "",
		Config:                 "",
		ServiceProviderAccount: "",
	}, nil
}

func (rc *MockReceptorClient) Discovered(ctx context.Context, in *receptor.Services, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: ""}, nil
}

func (rc *MockReceptorClient) Report(ctx context.Context, in *receptor.Finding, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: ""}, nil
}

func (rc *MockReceptorClient) Notify(ctx context.Context, in *receptor.JobResult, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
