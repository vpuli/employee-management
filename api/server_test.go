package api

import (
	"bytes"
	"employees/internal/db/mocks"
	"employees/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestServer(t *testing.T) (*Server, *mocks.Database) {
	logger := zap.NewNop()
	mockDB := mocks.NewDatabase(t)
	server := NewServer(":8080", logger, mockDB)
	return server, mockDB
}

func TestHandleCreateEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name       string
		employee   models.Employee
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name: "Valid Employee",
			employee: models.Employee{
				ID:        "1",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			setupMock: func(db *mocks.Database) {
				db.On("CreateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Database Error",
			employee: models.Employee{
				ID: "2",
			},
			setupMock: func(db *mocks.Database) {
				db.On("CreateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(errors.New("database error"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mockDB)

			// Create request
			payload, err := json.Marshal(tt.employee)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/employee", bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()

			// Handle request
			server.handleCreateEmployee(rr, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleGetEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	testEmp := &models.Employee{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	tests := []struct {
		name       string
		employeeID string
		setupMock  func(*mocks.Database)
		wantStatus int
		wantBody   *models.Employee
	}{
		{
			name:       "Existing Employee",
			employeeID: "1",
			setupMock: func(db *mocks.Database) {
				db.On("GetEmployee", "1").Return(testEmp, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   testEmp,
		},
		{
			name:       "Non-existent Employee",
			employeeID: "999",
			setupMock: func(db *mocks.Database) {
				db.On("GetEmployee", "999").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mockDB)

			// Create request
			req := httptest.NewRequest("GET", "/employee/"+tt.employeeID, nil)
			rr := httptest.NewRecorder()

			// Add URL parameters
			vars := map[string]string{"id": tt.employeeID}
			req = mux.SetURLVars(req, vars)

			// Handle request
			server.handleGetEmployee(rr, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != nil {
				var got models.Employee
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleUpdateEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	testEmp := &models.Employee{
		ID:        "1",
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
	}

	tests := []struct {
		name       string
		employeeID string
		update     models.Employee
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name:       "Valid Update",
			employeeID: "1",
			update:     *testEmp,
			setupMock: func(db *mocks.Database) {
				db.On("UpdateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Non-existent Employee",
			employeeID: "999",
			update:     *testEmp,
			setupMock: func(db *mocks.Database) {
				db.On("UpdateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mockDB)

			// Create request
			payload, err := json.Marshal(tt.update)
			require.NoError(t, err)

			req := httptest.NewRequest("PUT", "/employee/"+tt.employeeID, bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()

			// Add URL parameters
			vars := map[string]string{"id": tt.employeeID}
			req = mux.SetURLVars(req, vars)

			// Handle request
			server.handleUpdateEmployee(rr, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleDeleteEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name       string
		employeeID string
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name:       "Existing Employee",
			employeeID: "1",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteEmployee", "1").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Non-existent Employee",
			employeeID: "999",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteEmployee", "999").Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock(mockDB)

			// Create request
			req := httptest.NewRequest("DELETE", "/employee/"+tt.employeeID, nil)
			rr := httptest.NewRecorder()

			// Add URL parameters
			vars := map[string]string{"id": tt.employeeID}
			req = mux.SetURLVars(req, vars)

			// Handle request
			server.handleDeleteEmployee(rr, req)

			// Assert response
			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}
