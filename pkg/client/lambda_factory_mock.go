// Code generated by MockGen. DO NOT EDIT.
// Source: ./lambda_factory.go

// Package client is a generated GoMock package.
package client

import (
	reflect "reflect"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	gomock "github.com/golang/mock/gomock"
)

// MockLambdaCreator is a mock of LambdaCreator interface.
type MockLambdaCreator struct {
	ctrl     *gomock.Controller
	recorder *MockLambdaCreatorMockRecorder
}

// MockLambdaCreatorMockRecorder is the mock recorder for MockLambdaCreator.
type MockLambdaCreatorMockRecorder struct {
	mock *MockLambdaCreator
}

// NewMockLambdaCreator creates a new mock instance.
func NewMockLambdaCreator(ctrl *gomock.Controller) *MockLambdaCreator {
	mock := &MockLambdaCreator{ctrl: ctrl}
	mock.recorder = &MockLambdaCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLambdaCreator) EXPECT() *MockLambdaCreatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLambdaCreator) Create(config aws.Config) LambdaClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(LambdaClient)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockLambdaCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLambdaCreator)(nil).Create), config)
}
