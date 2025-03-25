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
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
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
			tt.setupMock(mockDB)

			payload, err := json.Marshal(tt.employee)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/employee", bytes.NewBuffer(payload))
			req.Header.Set("Authorization", "Bearer valid-token")
			rr := httptest.NewRecorder()

			server.handleCreateEmployee(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleGetEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	testEmp := &models.Employee{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	tests := []struct {
		name       string
		id         string
		setupMock  func(*mocks.Database)
		wantStatus int
		wantBody   *models.Employee
	}{
		{
			name: "Existing Employee",
			id:   "1",
			setupMock: func(db *mocks.Database) {
				db.On("GetEmployee", "1").Return(testEmp, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   testEmp,
		},
		{
			name: "Non-existent Employee",
			id:   "999",
			setupMock: func(db *mocks.Database) {
				db.On("GetEmployee", "999").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
		{
			name:       "Missing ID",
			id:         "",
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			req := httptest.NewRequest("GET", "/employee?id="+tt.id, nil)
			req.Header.Set("Authorization", "Bearer valid-token")
			rr := httptest.NewRecorder()

			server.handleGetEmployee(rr, req)

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

	tests := []struct {
		name       string
		id         string
		update     models.Employee
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name: "Valid Update",
			id:   "1",
			update: models.Employee{
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@example.com",
			},
			setupMock: func(db *mocks.Database) {
				db.On("UpdateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Non-existent Employee",
			id:   "999",
			update: models.Employee{
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@example.com",
			},
			setupMock: func(db *mocks.Database) {
				db.On("UpdateEmployee", mock.AnythingOfType("*models.Employee")).
					Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Missing ID",
			id:         "",
			update:     models.Employee{},
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			payload, err := json.Marshal(tt.update)
			require.NoError(t, err)

			req := httptest.NewRequest("PUT", "/employee?id="+tt.id, bytes.NewBuffer(payload))
			req.Header.Set("Authorization", "Bearer valid-token")
			rr := httptest.NewRecorder()

			server.handleUpdateEmployee(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleDeleteEmployee(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name       string
		id         string
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name: "Existing Employee",
			id:   "1",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteEmployee", "1").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Non-existent Employee",
			id:   "999",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteEmployee", "999").Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Missing ID",
			id:         "",
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			req := httptest.NewRequest("DELETE", "/employee?id="+tt.id, nil)
			req.Header.Set("Authorization", "Bearer valid-token")
			rr := httptest.NewRecorder()

			server.handleDeleteEmployee(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleCreateAdmin(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name       string
		admin      models.Admin
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name: "Valid Admin",
			admin: models.Admin{
				Email:    "admin@example.com",
				Password: "password123",
			},
			setupMock: func(db *mocks.Database) {
				db.On("CreateAdmin", mock.AnythingOfType("*models.Admin")).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing Email",
			admin: models.Admin{
				Password: "password123",
			},
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Password",
			admin: models.Admin{
				Email: "admin@example.com",
			},
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Database Error",
			admin: models.Admin{
				Email:    "admin@example.com",
				Password: "password123",
			},
			setupMock: func(db *mocks.Database) {
				db.On("CreateAdmin", mock.AnythingOfType("*models.Admin")).
					Return(errors.New("database error"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			payload, err := json.Marshal(tt.admin)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/admin", bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()

			server.handleCreateAdmin(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleGetAdmin(t *testing.T) {
	server, mockDB := setupTestServer(t)

	testAdmin := &models.Admin{
		ID:    1,
		Email: "admin@example.com",
	}

	tests := []struct {
		name       string
		email      string
		setupMock  func(*mocks.Database)
		wantStatus int
		wantBody   *models.Admin
	}{
		{
			name:  "Existing Admin",
			email: "admin@example.com",
			setupMock: func(db *mocks.Database) {
				db.On("GetAdmin", "admin@example.com").Return(testAdmin, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   testAdmin,
		},
		{
			name:  "Non-existent Admin",
			email: "nonexistent@example.com",
			setupMock: func(db *mocks.Database) {
				db.On("GetAdmin", "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
		{
			name:       "Missing Email",
			email:      "",
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			req := httptest.NewRequest("GET", "/admin?email="+tt.email, nil)
			rr := httptest.NewRecorder()

			server.handleGetAdmin(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != nil {
				var got models.Admin
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleUpdateAdmin(t *testing.T) {
	server, mockDB := setupTestServer(t)

	testAdmin := &models.Admin{
		ID:    1,
		Email: "admin@example.com",
	}

	tests := []struct {
		name   string
		email  string
		update struct {
			Password string `json:"password"`
		}
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name:  "Valid Update",
			email: "admin@example.com",
			update: struct {
				Password string `json:"password"`
			}{
				Password: "newpassword123",
			},
			setupMock: func(db *mocks.Database) {
				db.On("GetAdminByEmail", "admin@example.com").Return(testAdmin, nil)
				db.On("UpdateAdmin", mock.AnythingOfType("*models.Admin")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "Non-existent Admin",
			email: "nonexistent@example.com",
			update: struct {
				Password string `json:"password"`
			}{
				Password: "newpassword123",
			},
			setupMock: func(db *mocks.Database) {
				db.On("GetAdminByEmail", "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "Missing Email",
			email: "",
			update: struct {
				Password string `json:"password"`
			}{},
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			payload, err := json.Marshal(tt.update)
			require.NoError(t, err)

			req := httptest.NewRequest("PUT", "/admin?email="+tt.email, bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()

			server.handleUpdateAdmin(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestHandleDeleteAdmin(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name       string
		email      string
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name:  "Existing Admin",
			email: "admin@example.com",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteAdmin", "admin@example.com").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:  "Non-existent Admin",
			email: "nonexistent@example.com",
			setupMock: func(db *mocks.Database) {
				db.On("DeleteAdmin", "nonexistent@example.com").Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Missing Email",
			email:      "",
			setupMock:  func(db *mocks.Database) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			req := httptest.NewRequest("DELETE", "/admin?email="+tt.email, nil)
			rr := httptest.NewRecorder()

			server.handleDeleteAdmin(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	server, mockDB := setupTestServer(t)

	tests := []struct {
		name      string
		loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		setupMock  func(*mocks.Database)
		wantStatus int
	}{
		{
			name: "Valid Login",
			loginData: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "admin@example.com",
				Password: "password123",
			},
			setupMock: func(db *mocks.Database) {
				admin := &models.Admin{
					ID:       1,
					Email:    "admin@example.com",
					Password: "$2a$10$yourhashedpassword", // You'll need to set this up properly
				}
				db.On("GetAdmin", "admin@example.com").Return(admin, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid Credentials",
			loginData: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "admin@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(db *mocks.Database) {
				admin := &models.Admin{
					ID:       1,
					Email:    "admin@example.com",
					Password: "$2a$10$yourhashedpassword",
				}
				db.On("GetAdmin", "admin@example.com").Return(admin, nil)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Non-existent Admin",
			loginData: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(db *mocks.Database) {
				db.On("GetAdmin", "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockDB)

			payload, err := json.Marshal(tt.loginData)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(payload))
			rr := httptest.NewRecorder()

			server.LogIn(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantStatus == http.StatusOK {
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Contains(t, response, "token")
			}
			mockDB.AssertExpectations(t)
		})
	}
}
