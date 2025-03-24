package mock

import (
	"employees/internal/models"
	"testing"
)

func TestMockDB_CRUD(t *testing.T) {
	db := NewMockDB()

	// Test employee
	emp := &models.Employee{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	// Test Create
	t.Run("Create Employee", func(t *testing.T) {
		if err := db.CreateEmployee(emp); err != nil {
			t.Errorf("CreateEmployee() error = %v", err)
		}

		// Try creating duplicate
		if err := db.CreateEmployee(emp); err == nil {
			t.Error("CreateEmployee() expected error for duplicate employee")
		}
	})

	// Test Get
	t.Run("Get Employee", func(t *testing.T) {
		got, err := db.GetEmployee(emp.ID)
		if err != nil {
			t.Errorf("GetEmployee() error = %v", err)
		}
		if got.ID != emp.ID {
			t.Errorf("GetEmployee() = %v, want %v", got, emp)
		}

		// Try getting non-existent employee
		if _, err := db.GetEmployee("999"); err == nil {
			t.Error("GetEmployee() expected error for non-existent employee")
		}
	})

	// Test Update
	t.Run("Update Employee", func(t *testing.T) {
		emp.FirstName = "Jane"
		if err := db.UpdateEmployee(emp); err != nil {
			t.Errorf("UpdateEmployee() error = %v", err)
		}

		got, _ := db.GetEmployee(emp.ID)
		if got.FirstName != "Jane" {
			t.Errorf("UpdateEmployee() failed to update name, got = %v", got.FirstName)
		}

		// Try updating non-existent employee
		nonExistent := &models.Employee{ID: "999"}
		if err := db.UpdateEmployee(nonExistent); err == nil {
			t.Error("UpdateEmployee() expected error for non-existent employee")
		}
	})

	// Test Delete
	t.Run("Delete Employee", func(t *testing.T) {
		if err := db.DeleteEmployee(emp.ID); err != nil {
			t.Errorf("DeleteEmployee() error = %v", err)
		}

		// Verify deletion
		if _, err := db.GetEmployee(emp.ID); err == nil {
			t.Error("DeleteEmployee() failed to delete employee")
		}

		// Try deleting non-existent employee
		if err := db.DeleteEmployee("999"); err == nil {
			t.Error("DeleteEmployee() expected error for non-existent employee")
		}
	})
}
