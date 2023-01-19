// Code generated by MockGen. DO NOT EDIT.
// Source: ./lambda.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	types "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	gomock "github.com/golang/mock/gomock"
)

// MockLambdaClient is a mock of LambdaClient interface.
type MockLambdaClient struct {
	ctrl     *gomock.Controller
	recorder *MockLambdaClientMockRecorder
}

// MockLambdaClientMockRecorder is the mock recorder for MockLambdaClient.
type MockLambdaClientMockRecorder struct {
	mock *MockLambdaClient
}

// NewMockLambdaClient creates a new mock instance.
func NewMockLambdaClient(ctrl *gomock.Controller) *MockLambdaClient {
	mock := &MockLambdaClient{ctrl: ctrl}
	mock.recorder = &MockLambdaClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLambdaClient) EXPECT() *MockLambdaClientMockRecorder {
	return m.recorder
}

// ListFunctions mocks base method.
func (m *MockLambdaClient) ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFunctions", ctx)
	ret0, _ := ret[0].([]types.FunctionConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFunctions indicates an expected call of ListFunctions.
func (mr *MockLambdaClientMockRecorder) ListFunctions(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFunctions", reflect.TypeOf((*MockLambdaClient)(nil).ListFunctions), ctx)
}

// ListFunctionsWithRegion mocks base method.
func (m *MockLambdaClient) ListFunctionsWithRegion(ctx context.Context, region string) ([]types.FunctionConfiguration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFunctionsWithRegion", ctx, region)
	ret0, _ := ret[0].([]types.FunctionConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFunctionsWithRegion indicates an expected call of ListFunctionsWithRegion.
func (mr *MockLambdaClientMockRecorder) ListFunctionsWithRegion(ctx, region interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFunctionsWithRegion", reflect.TypeOf((*MockLambdaClient)(nil).ListFunctionsWithRegion), ctx, region)
}

// ListRuntimeValues mocks base method.
func (m *MockLambdaClient) ListRuntimeValues() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRuntimeValues")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListRuntimeValues indicates an expected call of ListRuntimeValues.
func (mr *MockLambdaClientMockRecorder) ListRuntimeValues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRuntimeValues", reflect.TypeOf((*MockLambdaClient)(nil).ListRuntimeValues))
}
