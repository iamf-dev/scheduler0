// Code generated by mockery v2.26.1. DO NOT EDIT.

package shared_repo

import (
	db "scheduler0/db"
	models "scheduler0/models"

	mock "github.com/stretchr/testify/mock"
)

// MockSharedRepo is an autogenerated mock type for the SharedRepo type
type MockSharedRepo struct {
	mock.Mock
}

// DeleteAsyncTasksLogs provides a mock function with given fields: _a0, committed, asyncTasks
func (_m *MockSharedRepo) DeleteAsyncTasksLogs(_a0 db.DataStore, committed bool, asyncTasks []models.AsyncTask) error {
	ret := _m.Called(_a0, committed, asyncTasks)

	var r0 error
	if rf, ok := ret.Get(0).(func(db.DataStore, bool, []models.AsyncTask) error); ok {
		r0 = rf(_a0, committed, asyncTasks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteExecutionLogs provides a mock function with given fields: _a0, committed, jobExecutionLogs
func (_m *MockSharedRepo) DeleteExecutionLogs(_a0 db.DataStore, committed bool, jobExecutionLogs []models.JobExecutionLog) error {
	ret := _m.Called(_a0, committed, jobExecutionLogs)

	var r0 error
	if rf, ok := ret.Get(0).(func(db.DataStore, bool, []models.JobExecutionLog) error); ok {
		r0 = rf(_a0, committed, jobExecutionLogs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetExecutionLogs provides a mock function with given fields: _a0, committed
func (_m *MockSharedRepo) GetExecutionLogs(_a0 db.DataStore, committed bool) ([]models.JobExecutionLog, error) {
	ret := _m.Called(_a0, committed)

	var r0 []models.JobExecutionLog
	var r1 error
	if rf, ok := ret.Get(0).(func(db.DataStore, bool) ([]models.JobExecutionLog, error)); ok {
		return rf(_a0, committed)
	}
	if rf, ok := ret.Get(0).(func(db.DataStore, bool) []models.JobExecutionLog); ok {
		r0 = rf(_a0, committed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.JobExecutionLog)
		}
	}

	if rf, ok := ret.Get(1).(func(db.DataStore, bool) error); ok {
		r1 = rf(_a0, committed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertAsyncTasksLogs provides a mock function with given fields: _a0, committed, asyncTasks
func (_m *MockSharedRepo) InsertAsyncTasksLogs(_a0 db.DataStore, committed bool, asyncTasks []models.AsyncTask) error {
	ret := _m.Called(_a0, committed, asyncTasks)

	var r0 error
	if rf, ok := ret.Get(0).(func(db.DataStore, bool, []models.AsyncTask) error); ok {
		r0 = rf(_a0, committed, asyncTasks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertExecutionLogs provides a mock function with given fields: _a0, committed, jobExecutionLogs
func (_m *MockSharedRepo) InsertExecutionLogs(_a0 db.DataStore, committed bool, jobExecutionLogs []models.JobExecutionLog) error {
	ret := _m.Called(_a0, committed, jobExecutionLogs)

	var r0 error
	if rf, ok := ret.Get(0).(func(db.DataStore, bool, []models.JobExecutionLog) error); ok {
		r0 = rf(_a0, committed, jobExecutionLogs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockSharedRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSharedRepo creates a new instance of MockSharedRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSharedRepo(t mockConstructorTestingTNewMockSharedRepo) *MockSharedRepo {
	mock := &MockSharedRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
