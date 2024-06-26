// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mock/mock.go service
//
// Package mock_service is a generated GoMock package.
package mock_service

import (
	url "net/url"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockVideo is a mock of Video interface.
type MockVideo struct {
	ctrl     *gomock.Controller
	recorder *MockVideoMockRecorder
}

// MockVideoMockRecorder is the mock recorder for MockVideo.
type MockVideoMockRecorder struct {
	mock *MockVideo
}

// NewMockVideo creates a new mock instance.
func NewMockVideo(ctrl *gomock.Controller) *MockVideo {
	mock := &MockVideo{ctrl: ctrl}
	mock.recorder = &MockVideoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVideo) EXPECT() *MockVideoMockRecorder {
	return m.recorder
}

// GenerateCDNUrl mocks base method.
func (m *MockVideo) GenerateCDNUrl(originalURL url.URL, clusterName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCDNUrl", originalURL, clusterName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCDNUrl indicates an expected call of GenerateCDNUrl.
func (mr *MockVideoMockRecorder) GenerateCDNUrl(originalURL, clusterName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCDNUrl", reflect.TypeOf((*MockVideo)(nil).GenerateCDNUrl), originalURL, clusterName)
}

// ValidateOriginalURL mocks base method.
func (m *MockVideo) ValidateOriginalURL(rawOriginalURL string) (url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateOriginalURL", rawOriginalURL)
	ret0, _ := ret[0].(url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateOriginalURL indicates an expected call of ValidateOriginalURL.
func (mr *MockVideoMockRecorder) ValidateOriginalURL(rawOriginalURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateOriginalURL", reflect.TypeOf((*MockVideo)(nil).ValidateOriginalURL), rawOriginalURL)
}
