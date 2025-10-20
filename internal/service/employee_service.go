package service

import (
	"context"
	"errors"

	"emplopyee-app-go/internal/dao"
	"emplopyee-app-go/internal/model"
)

var ErrNotFound = errors.New("employee not found")

type EmployeeService interface {
	CreateEmployee(ctx context.Context, in *model.Employee) (*model.Employee, error)
	UpdateEmployee(ctx context.Context, in *model.Employee) (*model.Employee, error)
	GetEmployee(ctx context.Context, id int64) (*model.Employee, error)
	ListEmployees(ctx context.Context) ([]*model.Employee, error)
	DeleteEmployee(ctx context.Context, id int64) error
}

// EmployeeService defines business methods for managing employees.
//
// Provides CRUD operations on Employee objects, allowing higher-level validation and logic
// beyond simple database access.
//
// Methods:
//   - CreateEmployee: Adds a new employee to the system.
//   - UpdateEmployee: Updates an existing employee (must exist).
//   - GetEmployee:    Fetches an employee by unique ID.
//   - ListEmployees:  Returns all employees, ordered by ID descending.
//   - DeleteEmployee: Removes an employee by ID (must exist).
//
// Implementation wraps an EmployeeDAO, surfacing application-level errors such as ErrNotFound.

// ErrNotFound is returned when no employee is found for a given query.

// NewEmployeeService returns a new EmployeeService with the provided EmployeeDAO;
// this enables separation of business logic from data access logic.

type employeeService struct {
	dao dao.EmployeeDAO
}

func NewEmployeeService(d dao.EmployeeDAO) EmployeeService {
	return &employeeService{dao: d}
}

func (s *employeeService) CreateEmployee(ctx context.Context, in *model.Employee) (*model.Employee, error) {
	return s.dao.Create(ctx, in)
}

func (s *employeeService) UpdateEmployee(ctx context.Context, in *model.Employee) (*model.Employee, error) {
	// check exists
	_, err := s.dao.GetByID(ctx, in.ID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.dao.Update(ctx, in)
}

func (s *employeeService) GetEmployee(ctx context.Context, id int64) (*model.Employee, error) {
	e, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *employeeService) ListEmployees(ctx context.Context) ([]*model.Employee, error) {
	return s.dao.GetAll(ctx)
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id int64) error {
	// verify exists
	_, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	return s.dao.Delete(ctx, id)
}
