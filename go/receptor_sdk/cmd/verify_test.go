package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trustero/api/go/receptor_sdk/client"
	. "github.com/trustero/api/go/receptor_sdk/client/mocks"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_sdk/cmd/mocks"
	"github.com/trustero/api/go/receptor_sdk/config"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func init() {
	config.InitLog("debug", "")
}
func TestVerifyWithTestCredentials(t *testing.T) {

	mockConfig := &mocks.Config{}
	cmd.Config = mockConfig
	mockConfig.On("CredentialsFromFlags").Return(&map[string]string{})
	mockConfig.On("Verify", mock.Anything).Return(true, nil)

	cmd.NoSave = true
	defer func() {
		cmd.NoSave = false
	}()
	err := cmd.Verify(nil, nil)

	assert.Nil(t, err)
	mockConfig.AssertExpectations(t)
}

func TestVerify(t *testing.T) {
	rawServiceCredentials := "{\"username\":\"test\",\"password\":\"test\"}"
	serviceCredentials := &struct{}{}
	token := "xyz-token-qwe"
	receptorClient := &Receptor{}
	receptorClient.On("Verified", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

	mockConfig := &mocks.Config{}
	cmd.Config = mockConfig
	mockConfig.On("UnmarshallCredentials", rawServiceCredentials).Return(serviceCredentials, nil)
	mockConfig.On("CredentialsFromFlags").Return(nil)
	mockConfig.On("Verify", serviceCredentials).Return(true, nil)
	mockConfig.On("ReceptorModelId").Return("trr-test", nil)

	mockFactory := &MockScopedClientFactory{
		Client: receptorClient,
		ReceptorConfiguration: &receptor_v1.ReceptorConfiguration{
			Config:                 "abc",
			Credential:             rawServiceCredentials,
			ReceptorObjectId:       "trr-test",
			ServiceProviderAccount: "SPA-abc",
		},
	}
	mockFactory.On("AuthScope", mock.Anything, mock.Anything).Return(nil)
	client.Factory = mockFactory

	err := cmd.Verify(nil, []string{token})

	assert.Nil(t, err)
	receptorClient.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}
