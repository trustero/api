package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/trustero/api/go/receptor_sdk/client"
	"github.com/trustero/api/go/receptor_v1"
)

type MockScopedClientFactory struct {
	Client                receptor_v1.ReceptorClient
	ReceptorConfiguration *receptor_v1.ReceptorConfiguration
	Context               context.Context
	mock.Mock
}

func (t *MockScopedClientFactory) AuthScope(modelId string, onAuth client.AuthClientCallback) (err error) {
	t.Called(modelId, onAuth)
	return onAuth(t.Context, t.Client, t.ReceptorConfiguration)
}
