package dao

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"emplopyee-app-go/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEmployeeDAO_GetByID_Success(t *testing.T) {
	// Arrange: Set up mocked SQL DB and test data

	db, mock, err := sqlmock.New() // Create a new SQL mock database and mock handler
	if err != nil {                // Check for error during mock DB creation
		t.Fatalf("sqlmock.New: %v", err) // Fail test if DB could not be created
	}
	defer db.Close() // Ensure DB is closed after test completes

	da := NewEmployeeDAO(db) // Instantiate EmployeeDAO with the mock DB

	cols := []string{"id", "first_name", "last_name", "email", "position", "created_at", "updated_at"} // Define the employee table columns
	now := time.Now().UTC()                                                                            // Fix a timestamp for created_at/updated_at fields
	row := sqlmock.NewRows(cols).                                                                      // Build a mock SQL row with given columns
														AddRow(int64(1), "Alice", "Smith", "alice@example.com", "Engineer", now, now) // Add a single employee test record

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employees WHERE id = ?")). // Expect a specific query to be executed against mock DB
											WithArgs(int64(1)). // Expect the argument value to be 1
											WillReturnRows(row) // Mock the query to return the prepared row

	// Act: Invoke the function under test

	got, err := da.GetByID(context.Background(), 1) // Attempt to fetch employee with ID=1 from the DAO

	// Assert: Validate results and mock expectations

	if err != nil { // Confirm there was no error when fetching employee
		t.Fatalf("GetByID error: %v", err) // Fail test on error
	}
	if got == nil { // Ensure a result was returned, not nil
		t.Fatalf("expected employee, got nil") // Fail test if nil is returned
	}
	want := &model.Employee{ // Construct the expected Employee values (ignoring timestamps)
		ID:        1,
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		Position:  "Engineer",
	}
	if got.ID != want.ID || got.FirstName != want.FirstName || got.LastName != want.LastName || got.Email != want.Email || got.Position != want.Position {
		// Check for equality of core fields; log error if mismatch
		t.Errorf("unexpected employee: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil { // Confirm all SQL mock expectations were actually triggered
		t.Errorf("unmet expectations: %v", err) // Report if not all mock behaviors were satisfied
	}
}

// TestEmployeeDAO_GetByID_Success tests the successful case for fetching an employee by ID.
// Lines below describe what is done at each stage:

// 1. db, mock, err := sqlmock.New()
// Creates a new in-memory mock SQL database and mock object for interaction.

// 2. defer db.Close()
// Schedules the closing of the mock database connection after test ends.

// 3. da := NewEmployeeDAO(db)
// Constructs a new EmployeeDAO using the mocked database.

// 4. cols := []string{"id", "first_name", "last_name", "email", "position", "created_at", "updated_at"}
// Defines all columns to be returned in the row.

// 5. now := time.Now().UTC()
// Sets up a common timestamp for both created_at and updated_at.

// 6. row := sqlmock.NewRows(cols).AddRow(...)
// Prepares a single row of employee test data to be returned by the mocked query result.

// 7. mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employees WHERE id = ?")).
//       WithArgs(int64(1)).
//       WillReturnRows(row)
// Configures the mock to expect the specific SELECT query, with argument value 1, returning the test row.

// 8. got, err := da.GetByID(context.Background(), 1)
// Calls GetByID to fetch employee with ID 1.

// 9. if err != nil { t.Fatalf("GetByID error: %v", err) }
// Asserts the fetch did not return an error.

// 10. if got == nil { t.Fatalf("expected employee, got nil") }
// Asserts the result should not be nil.

// 11. want := &model.Employee{ (fields...) }
// Defines the expected core fields of the employee object.

// 12. if got.ID != want.ID || ... { t.Errorf("unexpected employee: %+v", got) }
// Checks all non-time fields for exact match.

// 13. if err := mock.ExpectationsWereMet(); err != nil { ... }
// Verifies all mock expectations (queries, args) were satisfied during test execution.

func TestEmployeeDAO_GetByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	da := NewEmployeeDAO(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employees WHERE id = ?")).
		WithArgs(int64(42)).
		WillReturnError(sql.ErrNoRows)

	// Act
	got, err := da.GetByID(context.Background(), 42)

	// Assert
	if err == nil {
		t.Fatalf("expected error, got nil (employee: %+v)", got)
	}
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
