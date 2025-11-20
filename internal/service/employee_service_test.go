package service

import (
	"context"
	"errors"
	"testing"

	"emplopyee-app-go/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeDAO is a mock implementation of dao.EmployeeDAO
type MockEmployeeDAO struct {
	mock.Mock
}

func (m *MockEmployeeDAO) Create(ctx context.Context, e *model.Employee) (*model.Employee, error) {
	args := m.Called(ctx, e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeDAO) Update(ctx context.Context, e *model.Employee) (*model.Employee, error) {
	args := m.Called(ctx, e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeDAO) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeDAO) GetAll(ctx context.Context) ([]*model.Employee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Employee), args.Error(1)
}

func (m *MockEmployeeDAO) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateEmployee(t *testing.T) {
	mockDAO := new(MockEmployeeDAO)
	svc := NewEmployeeService(mockDAO)
	ctx := context.Background()

	input := &model.Employee{FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	expected := &model.Employee{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com"}

	mockDAO.On("Create", ctx, input).Return(expected, nil)

	result, err := svc.CreateEmployee(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockDAO.AssertExpectations(t)
}

func TestGetEmployee(t *testing.T) {
	ctx := context.Background()

	t.Run("Found", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		expected := &model.Employee{ID: 1, FirstName: "John"}
		mockDAO.On("GetByID", ctx, int64(1)).Return(expected, nil)

		result, err := svc.GetEmployee(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockDAO.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		mockDAO.On("GetByID", ctx, int64(999)).Return(nil, errors.New("db error"))

		result, err := svc.GetEmployee(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, result)
		mockDAO.AssertExpectations(t)
	})
}

func TestUpdateEmployee(t *testing.T) {
	ctx := context.Background()
	input := &model.Employee{ID: 1, FirstName: "Jane"}

	t.Run("Success", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		// First it checks if exists
		mockDAO.On("GetByID", ctx, int64(1)).Return(&model.Employee{ID: 1}, nil)
		// Then updates
		mockDAO.On("Update", ctx, input).Return(input, nil)

		result, err := svc.UpdateEmployee(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, input, result)
		mockDAO.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		mockDAO.On("GetByID", ctx, int64(1)).Return(nil, errors.New("not found"))

		result, err := svc.UpdateEmployee(ctx, input)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, result)
		mockDAO.AssertExpectations(t)
	})
}

func TestDeleteEmployee(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		mockDAO.On("GetByID", ctx, int64(1)).Return(&model.Employee{ID: 1}, nil)
		mockDAO.On("Delete", ctx, int64(1)).Return(nil)

		err := svc.DeleteEmployee(ctx, 1)
		assert.NoError(t, err)
		mockDAO.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDAO := new(MockEmployeeDAO)
		svc := NewEmployeeService(mockDAO)
		mockDAO.On("GetByID", ctx, int64(1)).Return(nil, errors.New("not found"))

		err := svc.DeleteEmployee(ctx, 1)
		assert.ErrorIs(t, err, ErrNotFound)
		mockDAO.AssertExpectations(t)
	})
}

func TestListEmployees(t *testing.T) {
	mockDAO := new(MockEmployeeDAO)
	svc := NewEmployeeService(mockDAO)
	ctx := context.Background()

	expected := []*model.Employee{
		{ID: 1, FirstName: "John"},
		{ID: 2, FirstName: "Jane"},
	}

	mockDAO.On("GetAll", ctx).Return(expected, nil)

	result, err := svc.ListEmployees(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockDAO.AssertExpectations(t)
}
