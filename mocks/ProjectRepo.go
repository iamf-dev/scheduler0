// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "scheduler0/models"

	mock "github.com/stretchr/testify/mock"

	utils "scheduler0/utils"
)

// ProjectRepo is an autogenerated mock type for the ProjectRepo type
type ProjectRepo struct {
	mock.Mock
}

// Count provides a mock function with given fields:
func (_m *ProjectRepo) Count() (uint64, *utils.GenericError) {
	ret := _m.Called()

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func() (uint64, *utils.GenericError)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() *utils.GenericError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// CreateOne provides a mock function with given fields: project
func (_m *ProjectRepo) CreateOne(project *models.ProjectModel) (uint64, *utils.GenericError) {
	ret := _m.Called(project)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(*models.ProjectModel) (uint64, *utils.GenericError)); ok {
		return rf(project)
	}
	if rf, ok := ret.Get(0).(func(*models.ProjectModel) uint64); ok {
		r0 = rf(project)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(*models.ProjectModel) *utils.GenericError); ok {
		r1 = rf(project)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: project
func (_m *ProjectRepo) DeleteOneByID(project models.ProjectModel) (uint64, *utils.GenericError) {
	ret := _m.Called(project)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.ProjectModel) (uint64, *utils.GenericError)); ok {
		return rf(project)
	}
	if rf, ok := ret.Get(0).(func(models.ProjectModel) uint64); ok {
		r0 = rf(project)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(models.ProjectModel) *utils.GenericError); ok {
		r1 = rf(project)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetBatchProjectsByIDs provides a mock function with given fields: projectIds
func (_m *ProjectRepo) GetBatchProjectsByIDs(projectIds []uint64) ([]models.ProjectModel, *utils.GenericError) {
	ret := _m.Called(projectIds)

	var r0 []models.ProjectModel
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func([]uint64) ([]models.ProjectModel, *utils.GenericError)); ok {
		return rf(projectIds)
	}
	if rf, ok := ret.Get(0).(func([]uint64) []models.ProjectModel); ok {
		r0 = rf(projectIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ProjectModel)
		}
	}

	if rf, ok := ret.Get(1).(func([]uint64) *utils.GenericError); ok {
		r1 = rf(projectIds)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// GetOneByID provides a mock function with given fields: project
func (_m *ProjectRepo) GetOneByID(project *models.ProjectModel) *utils.GenericError {
	ret := _m.Called(project)

	var r0 *utils.GenericError
	if rf, ok := ret.Get(0).(func(*models.ProjectModel) *utils.GenericError); ok {
		r0 = rf(project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.GenericError)
		}
	}

	return r0
}

// GetOneByName provides a mock function with given fields: project
func (_m *ProjectRepo) GetOneByName(project *models.ProjectModel) *utils.GenericError {
	ret := _m.Called(project)

	var r0 *utils.GenericError
	if rf, ok := ret.Get(0).(func(*models.ProjectModel) *utils.GenericError); ok {
		r0 = rf(project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.GenericError)
		}
	}

	return r0
}

// List provides a mock function with given fields: offset, limit
func (_m *ProjectRepo) List(offset uint64, limit uint64) ([]models.ProjectModel, *utils.GenericError) {
	ret := _m.Called(offset, limit)

	var r0 []models.ProjectModel
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(uint64, uint64) ([]models.ProjectModel, *utils.GenericError)); ok {
		return rf(offset, limit)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64) []models.ProjectModel); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ProjectModel)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64) *utils.GenericError); ok {
		r1 = rf(offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: project
func (_m *ProjectRepo) UpdateOneByID(project models.ProjectModel) (uint64, *utils.GenericError) {
	ret := _m.Called(project)

	var r0 uint64
	var r1 *utils.GenericError
	if rf, ok := ret.Get(0).(func(models.ProjectModel) (uint64, *utils.GenericError)); ok {
		return rf(project)
	}
	if rf, ok := ret.Get(0).(func(models.ProjectModel) uint64); ok {
		r0 = rf(project)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(models.ProjectModel) *utils.GenericError); ok {
		r1 = rf(project)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.GenericError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewProjectRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewProjectRepo creates a new instance of ProjectRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProjectRepo(t mockConstructorTestingTNewProjectRepo) *ProjectRepo {
	mock := &ProjectRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}