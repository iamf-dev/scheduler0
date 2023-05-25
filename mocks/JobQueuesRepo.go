// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "scheduler0/models"

	mock "github.com/stretchr/testify/mock"
)

// JobQueuesRepo is an autogenerated mock type for the JobQueuesRepo type
type JobQueuesRepo struct {
	mock.Mock
}

// GetLastJobQueueLogForNode provides a mock function with given fields: nodeId, version
func (_m *JobQueuesRepo) GetLastJobQueueLogForNode(nodeId uint64, version uint64) []models.JobQueueLog {
	ret := _m.Called(nodeId, version)

	var r0 []models.JobQueueLog
	if rf, ok := ret.Get(0).(func(uint64, uint64) []models.JobQueueLog); ok {
		r0 = rf(nodeId, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.JobQueueLog)
		}
	}

	return r0
}

// GetLastVersion provides a mock function with given fields:
func (_m *JobQueuesRepo) GetLastVersion() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

type mockConstructorTestingTNewJobQueuesRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewJobQueuesRepo creates a new instance of JobQueuesRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJobQueuesRepo(t mockConstructorTestingTNewJobQueuesRepo) *JobQueuesRepo {
	mock := &JobQueuesRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
