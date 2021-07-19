// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/tcsrdrserver/portfolio_response_struct.go

// Package tcsrdrserver is a generated GoMock package.
package tcsrdrserver

import (
	reflect "reflect"
	tcsrdrserver "tcsrdr/internal/tcsrdrserver"

	gomock "github.com/golang/mock/gomock"
)

// MockTCSPortfolioData is a mock of TCSPortfolioData interface.
type MockTCSPortfolioData struct {
	ctrl     *gomock.Controller
	recorder *MockTCSPortfolioDataMockRecorder
}

// MockTCSPortfolioDataMockRecorder is the mock recorder for MockTCSPortfolioData.
type MockTCSPortfolioDataMockRecorder struct {
	mock *MockTCSPortfolioData
}

// NewMockTCSPortfolioData creates a new mock instance.
func NewMockTCSPortfolioData(ctrl *gomock.Controller) *MockTCSPortfolioData {
	mock := &MockTCSPortfolioData{ctrl: ctrl}
	mock.recorder = &MockTCSPortfolioDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTCSPortfolioData) EXPECT() *MockTCSPortfolioDataMockRecorder {
	return m.recorder
}

// getData mocks base method.
func (m *MockTCSPortfolioData) getData(response *tcsrdrserver.TCSPortfolioResponse) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "getData", response)
}

// getData indicates an expected call of getData.
func (mr *MockTCSPortfolioDataMockRecorder) getData(response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getData", reflect.TypeOf((*MockTCSPortfolioData)(nil).getData), response)
}
