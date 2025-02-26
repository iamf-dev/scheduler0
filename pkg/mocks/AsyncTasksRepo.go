// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "scheduler0/pkg/models"

	mock "github.com/stretchr/testify/mock"

	utils "scheduler0/pkg/utils"
)

// AsyncTasksRepo is an autogenerated mock type for the AsyncTasksRepo type
type AsyncTasksRepo struct {
	mock.Mock
}

// BatchInsert provides a mock function with given fields: tasks, committed
func (_m *AsyncTasksRepo) BatchInsert(tasks []models.AsyncTask, committed bool) ([]uint64, *utils.GenericError) {
	ret := _m.Called(tasks, committed)

	var r0 []uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func([]models.AsyncTask, bool) ([]uint64, *utils.GenericError)); ok {
		return rf(tasks, committed)
	}
	if rf, ok := ret.Get(0).(func([]models.AsyncTask, bool) []uint64); ok {
		r0 = rf(tasks, committed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func([]models.AsyncTask, bool) *utils.GenericError); ok {
		r1 = rf(tasks, committed)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetAllTasks provides a mock function with given fields: committed
func (_m *AsyncTasksRepo) GetAllTasks(committed bool) ([]models.AsyncTask, *utils.GenericError) {
	ret := _m.Called(committed)

	var r0 []models.AsyncTask
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(bool) ([]models.AsyncTask, *utils.GenericError)); ok {
		return rf(committed)
	}
	if rf, ok := ret.Get(0).(func(bool) []models.AsyncTask); ok {
		r0 = rf(committed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.AsyncTask)
		}
	}

	if rf, ok := ret.Get(1).(func(bool) *utils.GenericError); ok {
		r1 = rf(committed)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetTask provides a mock function with given fields: taskId
func (_m *AsyncTasksRepo) GetTask(taskId uint64) (*models.AsyncTask, *utils.GenericError) {
	ret := _m.Called(taskId)

	var r0 *models.AsyncTask
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64) (*models.AsyncTask, *utils.GenericError)); ok {
		return rf(taskId)
	}
	if rf, ok := ret.Get(0).(func(uint64) *models.AsyncTask); ok {
		r0 = rf(taskId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AsyncTask)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) *utils.GenericError); ok {
		r1 = rf(taskId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// RaftBatchInsert provides a mock function with given fields: tasks
func (_m *AsyncTasksRepo) RaftBatchInsert(tasks []models.AsyncTask) ([]uint64, *utils.GenericError) {
	ret := _m.Called(tasks)

	var r0 []uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func([]models.AsyncTask) ([]uint64, *utils.GenericError)); ok {
		return rf(tasks)
	}
	if rf, ok := ret.Get(0).(func([]models.AsyncTask) []uint64); ok {
		r0 = rf(tasks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func([]models.AsyncTask) *utils.GenericError); ok {
		r1 = rf(tasks)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// RaftUpdateTaskState provides a mock function with given fields: task, state, output
func (_m *AsyncTasksRepo) RaftUpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError {
	ret := _m.Called(task, state, output)

	var r0 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.AsyncTask, models.AsyncTaskState, string) *utils.GenericError); ok {
		r0 = rf(task, state, output)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.GenericError)
		}
	}

	return r0
}

// UpdateTaskState provides a mock function with given fields: task, state, output
func (_m *AsyncTasksRepo) UpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError {
	ret := _m.Called(task, state, output)

	var r0 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.AsyncTask, models.AsyncTaskState, string) *utils.GenericError); ok {
		r0 = rf(task, state, output)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.GenericError)
		}
	}

	return r0
}

type mockConstructorTestingTNewAsyncTasksRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewAsyncTasksRepo creates a new instance of AsyncTasksRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAsyncTasksRepo(t mockConstructorTestingTNewAsyncTasksRepo) *AsyncTasksRepo {
	mock := &AsyncTasksRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
