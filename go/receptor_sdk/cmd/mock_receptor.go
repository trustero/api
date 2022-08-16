// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"context"
	"fmt"

	receptor "github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type mockReceptorClient struct{}

const header = "========\nReceptor."
const footer = "========\n\n"

// Verified implements a mock [receptor_v1.Receptor.Verified] method for testing.
func (rc *mockReceptorClient) Verified(ctx context.Context, in *receptor.Credential, opts ...grpc.CallOption) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	println(header + "Verified(...)")
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println(string(yamld))
	}
	println(footer)
	return
}

// Verified implements a mock [receptor_v1.Receptor.GetConfiguration] method for testing.
func (rc *mockReceptorClient) GetConfiguration(ctx context.Context, in *receptor.ReceptorOID, opts ...grpc.CallOption) (c *receptor.ReceptorConfiguration, err error) {
	c = &receptor.ReceptorConfiguration{
		ReceptorObjectId:       "",
		Credential:             "",
		Config:                 "",
		ServiceProviderAccount: "",
	}

	println(header + "GetConfiguration(...)")
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println(string(yamld))
	}
	println(footer)
	return
}

// Verified implements a mock [receptor_v1.Receptor.Discovered] method for testing.
func (rc *mockReceptorClient) Discovered(ctx context.Context, in *receptor.ServiceEntities, opts ...grpc.CallOption) (s *wrapperspb.StringValue, err error) {
	s = &wrapperspb.StringValue{Value: ""}

	println(header + "Discovered(...)")
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println(string(yamld))
	}
	println(footer)
	return
}

// Verified implements a mock [receptor_v1.Receptor.Report] method for testing.
func (rc *mockReceptorClient) Report(ctx context.Context, in *receptor.Finding, opts ...grpc.CallOption) (s *wrapperspb.StringValue, err error) {
	println(header + "Report(...)")

	println("Entities")
	var yamld string
	if yamld, err = toYaml(in.Entities); err == nil {
		println(string(yamld))
	}
	println()

	println("Evidences")
	for _, ev := range in.Evidences {

		t := ev.GetStruct()
		var headers []string
		var rows [][]string
		headers, rows, err = t.Tabulate()
		for _, header := range headers {
			fmt.Printf("%12s ", header)
		}
		println()

		for _, row := range rows {
			for _, col := range row {
				fmt.Printf("%12s ", col)
			}
			println()
		}
	}

	if err != nil {
		println(err)
	}

	println(footer)
	return
}

// Verified implements a mock [receptor_v1.Receptor.Notify] method for testing.
func (rc *mockReceptorClient) Notify(ctx context.Context, in *receptor.JobResult, opts ...grpc.CallOption) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}

	println(header + "Notify(...)")
	var yamld string
	if yamld, err = toYaml(in); err == nil {
		println(string(yamld))
	}
	println(footer)
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
