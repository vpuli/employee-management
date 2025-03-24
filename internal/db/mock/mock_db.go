package mock

import (
	"employees/internal/models"
	"fmt"
	"sync"
)

// MockDB implements the Database interface for testing
type MockDB struct {
	mu        sync.RWMutex
	employees map[string]*models.Employee
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		employees: make(map[string]*models.Employee),
	}
}

func (m *MockDB) CreateEmployee(emp *models.Employee) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.employees[emp.ID]; exists {
		return fmt.Errorf("employee with ID %s already exists", emp.ID)
	}

	m.employees[emp.ID] = emp
	return nil
}

func (m *MockDB) GetEmployee(id string) (*models.Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	emp, exists := m.employees[id]
	if !exists {
		return nil, fmt.Errorf("employee with ID %s not found", id)
	}

	return emp, nil
}

func (m *MockDB) UpdateEmployee(emp *models.Employee) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.employees[emp.ID]; !exists {
		return fmt.Errorf("employee with ID %s not found", emp.ID)
	}

	m.employees[emp.ID] = emp
	return nil
}

func (m *MockDB) DeleteEmployee(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.employees[id]; !exists {
		return fmt.Errorf("employee with ID %s not found", id)
	}

	delete(m.employees, id)
	return nil
}

func (m *MockDB) Close() error {
	return nil
}
