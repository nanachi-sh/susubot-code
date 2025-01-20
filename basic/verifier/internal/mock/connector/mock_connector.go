// Code generated by MockGen. DO NOT EDIT.
// Source: caller/connector/connector.go
//
// Generated by this command:
//
//	mockgen -source=caller/connector/connector.go -destination=mock/connector/mock_connector.go
//

// Package mock_connectorclient is a generated GoMock package.
package mock_connectorclient

import (
	context "context"
	reflect "reflect"

	connector "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/connector"
	connector0 "github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/connector"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockConnector is a mock of Connector interface.
type MockConnector struct {
	ctrl     *gomock.Controller
	recorder *MockConnectorMockRecorder
	isgomock struct{}
}

// MockConnectorMockRecorder is the mock recorder for MockConnector.
type MockConnectorMockRecorder struct {
	mock *MockConnector
}

// NewMockConnector creates a new mock instance.
func NewMockConnector(ctrl *gomock.Controller) *MockConnector {
	mock := &MockConnector{ctrl: ctrl}
	mock.recorder = &MockConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnector) EXPECT() *MockConnectorMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConnector) Close(ctx context.Context, in *connector.Empty, opts ...grpc.CallOption) (*connector.BasicResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Close", varargs...)
	ret0, _ := ret[0].(*connector.BasicResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Close indicates an expected call of Close.
func (mr *MockConnectorMockRecorder) Close(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConnector)(nil).Close), varargs...)
}

// Connect mocks base method.
func (m *MockConnector) Connect(ctx context.Context, in *connector.ConnectRequest, opts ...grpc.CallOption) (*connector.ConnectResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Connect", varargs...)
	ret0, _ := ret[0].(*connector.ConnectResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect.
func (mr *MockConnectorMockRecorder) Connect(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockConnector)(nil).Connect), varargs...)
}

// Read mocks base method.
func (m *MockConnector) Read(ctx context.Context, in *connector.Empty, opts ...grpc.CallOption) (connector0.Connector_ReadClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Read", varargs...)
	ret0, _ := ret[0].(connector0.Connector_ReadClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockConnectorMockRecorder) Read(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConnector)(nil).Read), varargs...)
}

// Write mocks base method.
func (m *MockConnector) Write(ctx context.Context, in *connector.WriteRequest, opts ...grpc.CallOption) (*connector.BasicResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Write", varargs...)
	ret0, _ := ret[0].(*connector.BasicResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockConnectorMockRecorder) Write(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockConnector)(nil).Write), varargs...)
}
