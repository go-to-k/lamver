// Code generated by MockGen. DO NOT EDIT.
// Source: ./ec2.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	gomock "github.com/golang/mock/gomock"
)

// MockEC2Client is a mock of EC2Client interface.
type MockEC2Client struct {
	ctrl     *gomock.Controller
	recorder *MockEC2ClientMockRecorder
}

// MockEC2ClientMockRecorder is the mock recorder for MockEC2Client.
type MockEC2ClientMockRecorder struct {
	mock *MockEC2Client
}

// NewMockEC2Client creates a new mock instance.
func NewMockEC2Client(ctrl *gomock.Controller) *MockEC2Client {
	mock := &MockEC2Client{ctrl: ctrl}
	mock.recorder = &MockEC2ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEC2Client) EXPECT() *MockEC2ClientMockRecorder {
	return m.recorder
}

// DescribeRegions mocks base method.
func (m *MockEC2Client) DescribeRegions(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeRegions", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeRegions indicates an expected call of DescribeRegions.
func (mr *MockEC2ClientMockRecorder) DescribeRegions(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeRegions", reflect.TypeOf((*MockEC2Client)(nil).DescribeRegions), ctx)
}

// MockEC2SDKClient is a mock of EC2SDKClient interface.
type MockEC2SDKClient struct {
	ctrl     *gomock.Controller
	recorder *MockEC2SDKClientMockRecorder
}

// MockEC2SDKClientMockRecorder is the mock recorder for MockEC2SDKClient.
type MockEC2SDKClientMockRecorder struct {
	mock *MockEC2SDKClient
}

// NewMockEC2SDKClient creates a new mock instance.
func NewMockEC2SDKClient(ctrl *gomock.Controller) *MockEC2SDKClient {
	mock := &MockEC2SDKClient{ctrl: ctrl}
	mock.recorder = &MockEC2SDKClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEC2SDKClient) EXPECT() *MockEC2SDKClientMockRecorder {
	return m.recorder
}

// DescribeRegions mocks base method.
func (m *MockEC2SDKClient) DescribeRegions(ctx context.Context, params *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeRegions", varargs...)
	ret0, _ := ret[0].(*ec2.DescribeRegionsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeRegions indicates an expected call of DescribeRegions.
func (mr *MockEC2SDKClientMockRecorder) DescribeRegions(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeRegions", reflect.TypeOf((*MockEC2SDKClient)(nil).DescribeRegions), varargs...)
}
