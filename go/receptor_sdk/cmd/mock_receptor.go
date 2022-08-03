// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"context"

	receptor "github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type MockReceptorClient struct{}

func (rc *MockReceptorClient) Verified(ctx context.Context, in *receptor.Credential, opts ...grpc.CallOption) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println("========")
		println("Verified credentials...")
		println(string(yamld))
		println()
	}
	return
}

func (rc *MockReceptorClient) GetConfiguration(ctx context.Context, in *receptor.ReceptorOID, opts ...grpc.CallOption) (c *receptor.ReceptorConfiguration, err error) {
	c = &receptor.ReceptorConfiguration{
		ReceptorObjectId:       "",
		Credential:             "",
		Config:                 "",
		ServiceProviderAccount: "",
	}
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println("========")
		println("GetConfiguration with receptor ID...")
		println(string(yamld))
		println()
	}
	return
}

func (rc *MockReceptorClient) Discovered(ctx context.Context, in *receptor.Services, opts ...grpc.CallOption) (s *wrapperspb.StringValue, err error) {
	s = &wrapperspb.StringValue{Value: ""}
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println("========")
		println("Discovered services...")
		println(string(yamld))
		println()
	}
	return
}

func (rc *MockReceptorClient) Report(ctx context.Context, in *receptor.Finding, opts ...grpc.CallOption) (s *wrapperspb.StringValue, err error) {
	s = &wrapperspb.StringValue{Value: ""}
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println("========")
		println("Report findings...")
		println(string(yamld))
		println()
	}
	return
}

func (rc *MockReceptorClient) Notify(ctx context.Context, in *receptor.JobResult, opts ...grpc.CallOption) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println("========")
		println("Notify job result...")
		println(string(yamld))
		println()
	}
	return
}

func toYaml(v interface{}) (yamld string, err error) {
	var bytes []byte
	if bytes, err = yaml.Marshal(v); err != nil {
		return
	}
	yamld = string(bytes)
	return
}
