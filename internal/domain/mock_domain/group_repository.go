// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/waliqueiroz/mystery-gifter-api/internal/domain (interfaces: GroupRepository)
//
// Generated by this command:
//
//	mockgen -destination mock_domain/group_repository.go . GroupRepository
//

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	domain "github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockGroupRepository is a mock of GroupRepository interface.
type MockGroupRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepositoryMockRecorder
	isgomock struct{}
}

// MockGroupRepositoryMockRecorder is the mock recorder for MockGroupRepository.
type MockGroupRepositoryMockRecorder struct {
	mock *MockGroupRepository
}

// NewMockGroupRepository creates a new mock instance.
func NewMockGroupRepository(ctrl *gomock.Controller) *MockGroupRepository {
	mock := &MockGroupRepository{ctrl: ctrl}
	mock.recorder = &MockGroupRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupRepository) EXPECT() *MockGroupRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockGroupRepository) Create(ctx context.Context, group domain.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockGroupRepositoryMockRecorder) Create(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGroupRepository)(nil).Create), ctx, group)
}

// GetByID mocks base method.
func (m *MockGroupRepository) GetByID(ctx context.Context, groupID string) (*domain.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, groupID)
	ret0, _ := ret[0].(*domain.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockGroupRepositoryMockRecorder) GetByID(ctx, groupID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockGroupRepository)(nil).GetByID), ctx, groupID)
}

// Update mocks base method.
func (m *MockGroupRepository) Update(ctx context.Context, group domain.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockGroupRepositoryMockRecorder) Update(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGroupRepository)(nil).Update), ctx, group)
}
