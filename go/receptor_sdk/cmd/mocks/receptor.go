package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/trustero/api/go/receptor_sdk/cmd"
)

type Config struct {
	mock.Mock
}

func (m *Config) Verify(serviceCredentials interface{}) (ok bool, err error) {
	args := m.Called(serviceCredentials)
	return args.Get(0).(bool), args.Error(1)
}

func (m *Config) Discover(serviceCredentials interface{}) (services []*cmd.Service, err error) {
	args := m.Called(serviceCredentials)
	return args.Get(0).([]*cmd.Service), args.Error(1)
}

// UnmarshallCredentials deserializes the credentials json string and returns the result as a struct pointer.
func (m *Config) UnmarshallCredentials(credentials string) (result interface{}, err error) {
	args := m.Called(credentials)
	return args.Get(0), args.Error(1)
}

func (m *Config) GetReporters() []cmd.Reporter {
	args := m.Called()
	return args.Get(0).([]cmd.Reporter)
}

func (m *Config) CredentialsFromFlags() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *Config) ReceptorModelId() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *Config) ServiceModelId() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *Config) EvidenceModelId() string {
	args := m.Called()
	return args.Get(0).(string)
}
