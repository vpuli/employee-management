// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "employees/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *Database) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAdmin provides a mock function with given fields: admin
func (_m *Database) CreateAdmin(admin *models.Admin) error {
	ret := _m.Called(admin)

	if len(ret) == 0 {
		panic("no return value specified for CreateAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Admin) error); ok {
		r0 = rf(admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateEmployee provides a mock function with given fields: emp
func (_m *Database) CreateEmployee(emp *models.Employee) error {
	ret := _m.Called(emp)

	if len(ret) == 0 {
		panic("no return value specified for CreateEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Employee) error); ok {
		r0 = rf(emp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAdmin provides a mock function with given fields: email
func (_m *Database) DeleteAdmin(email string) error {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEmployee provides a mock function with given fields: id
func (_m *Database) DeleteEmployee(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAdmin provides a mock function with given fields: email
func (_m *Database) GetAdmin(email string) (*models.Admin, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetAdmin")
	}

	var r0 *models.Admin
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Admin, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Admin); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Admin)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdminByEmail provides a mock function with given fields: email
func (_m *Database) GetAdminByEmail(email string) (*models.Admin, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetAdminByEmail")
	}

	var r0 *models.Admin
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Admin, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Admin); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Admin)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEmployee provides a mock function with given fields: id
func (_m *Database) GetEmployee(id string) (*models.Employee, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetEmployee")
	}

	var r0 *models.Employee
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Employee, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Employee); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Employee)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAdmin provides a mock function with given fields: admin
func (_m *Database) UpdateAdmin(admin *models.Admin) error {
	ret := _m.Called(admin)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Admin) error); ok {
		r0 = rf(admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateEmployee provides a mock function with given fields: emp
func (_m *Database) UpdateEmployee(emp *models.Employee) error {
	ret := _m.Called(emp)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Employee) error); ok {
		r0 = rf(emp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDatabase creates a new instance of Database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Database {
	mock := &Database{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
