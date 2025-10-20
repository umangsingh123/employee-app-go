package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"emplopyee-app-go/internal/model"

	"github.com/jmoiron/sqlx"
)

type EmployeeDAO interface {
	Create(ctx context.Context, e *model.Employee) (*model.Employee, error)
	Update(ctx context.Context, e *model.Employee) (*model.Employee, error)
	GetByID(ctx context.Context, id int64) (*model.Employee, error)
	GetAll(ctx context.Context) ([]*model.Employee, error)
	Delete(ctx context.Context, id int64) error
}

type employeeDAO struct {
	db *sqlx.DB
}

func NewEmployeeDAO(db *sql.DB) EmployeeDAO {
	return &employeeDAO{db: sqlx.NewDb(db, "sqlite3")}
}

// EmployeeDAO provides methods for CRUD operations on Employee model.
//
// Methods:
//   - Create: Inserts a new employee record into the database.
//   - Update: Updates an existing employee record.
//   - GetByID: Retrieves an employee by their unique ID.
//   - GetAll: Retrieves all employees from the database.
//   - Delete: Removes an employee record by ID.
//
// NewEmployeeDAO constructs a new EmployeeDAO backed by a sql.DB.
//
/*
Example usage:

	db, err := sql.Open("sqlite3", dsn)
	if err != nil { // handle error }
	dao := dao.NewEmployeeDAO(db)
	emp, err := dao.GetByID(ctx, 1)
*/

func (d *employeeDAO) Create(ctx context.Context, e *model.Employee) (*model.Employee, error) {
	query := `INSERT INTO employees (first_name, last_name, email, position, created_at, updated_at)
              VALUES (:first_name, :last_name, :email, :position, :created_at, :updated_at)`
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now
	res, err := d.db.NamedExecContext(ctx, query, e)
	if err != nil {
		return nil, fmt.Errorf("insert employee: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}
	e.ID = id
	return e, nil
}

func (d *employeeDAO) Update(ctx context.Context, e *model.Employee) (*model.Employee, error) {
	e.UpdatedAt = time.Now().UTC()
	query := `UPDATE employees SET first_name=:first_name, last_name=:last_name, email=:email, position=:position, updated_at=:updated_at WHERE id=:id`
	res, err := d.db.NamedExecContext(ctx, query, e)
	if err != nil {
		return nil, fmt.Errorf("update employee: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return e, nil
}

func (d *employeeDAO) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	var e model.Employee
	err := d.db.GetContext(ctx, &e, "SELECT * FROM employees WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (d *employeeDAO) GetAll(ctx context.Context) ([]*model.Employee, error) {
	var list []*model.Employee
	err := d.db.SelectContext(ctx, &list, "SELECT * FROM employees ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *employeeDAO) Delete(ctx context.Context, id int64) error {
	res, err := d.db.ExecContext(ctx, "DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
