package api

import (
	"employees/api/auth"
	"employees/api/middlewares"
	"employees/internal/db"
	"employees/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	logger     *zap.Logger
	router     *http.ServeMux
	listenAddr string
	db         db.Database
}

func NewServer(listenAddr string, logger *zap.Logger, db db.Database) *Server {
	return &Server{
		logger:     logger,
		listenAddr: listenAddr,
		router:     http.NewServeMux(),
		db:         db,
	}
}

func (s *Server) Start() error {
	s.router.HandleFunc("/employee", middlewares.SetMiddlewareAuthentication(s.handleEmployee))
	s.router.HandleFunc("/admin", s.handleAdmin)
	s.router.HandleFunc("/login", s.LogIn)
	return http.ListenAndServe(s.listenAddr, s.router)
}

func (s *Server) handleEmployee(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.handleCreateEmployee(w, r)
	case "GET":
		s.handleGetEmployee(w, r)
	case "PUT":
		s.handleUpdateEmployee(w, r)
	case "DELETE":
		s.handleDeleteEmployee(w, r)
	}
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.handleCreateAdmin(w, r)
	case "GET":
		s.handleGetAdmin(w, r)
	case "PUT":
		s.handleUpdateAdmin(w, r)
	case "DELETE":
		s.handleDeleteAdmin(w, r)
	}
}

func (s *Server) handleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.db.CreateEmployee(&emp); err != nil {
		s.logger.Error("Employee creation failed", zap.Error(err))
		http.Error(w, "Failed to create employee", http.StatusBadRequest)
		return
	}

	s.logger.Info("Employee created", zap.Any("employee", emp))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(emp); err != nil {
		s.logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (s *Server) handleGetEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		s.logger.Error("ID is required")
		return
	}

	emp, err := s.db.GetEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		s.logger.Error("Employee not found", zap.Error(err))
		return
	}

	s.logger.Info("Employee retrieved", zap.Any("employee", emp))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

func (s *Server) handleUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		s.logger.Error("ID is required")
		return
	}

	intId, ok := validateId(id)
	if !ok {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		s.logger.Error("Invalid ID")
		return
	}

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	emp.ID = intId // Ensure ID matches URL parameter
	err := s.db.UpdateEmployee(&emp)
	if err != nil {
		s.logger.Error("Employee update failed", zap.Error(err))
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	s.logger.Info("Employee updated", zap.Any("employee", emp))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(emp); err != nil {
		s.logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (s *Server) handleDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		s.logger.Error("ID is required")
		return
	}

	if err := s.db.DeleteEmployee(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		s.logger.Error("Employee deletion failed", zap.Error(err))
		return
	}

	s.logger.Info("Employee deleted", zap.Any("employeeId", id))

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if admin.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		s.logger.Error("Email is required")
		return
	}

	if admin.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		s.logger.Error("Password is required")
		return
	}

	if err := admin.BeforeSave(); err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		http.Error(w, "Failed to create admin", http.StatusInternalServerError)
		return
	}

	if err := s.db.CreateAdmin(&admin); err != nil {
		s.logger.Error("Admin creation failed", zap.Error(err))
		http.Error(w, "Failed to create admin", http.StatusBadRequest)
		return
	}

	s.logger.Info("Admin created", zap.Any("admin", admin))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(admin); err != nil {
		s.logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (s *Server) handleGetAdmin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		s.logger.Error("Email is required")
		return
	}

	admin, err := s.db.GetAdmin(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		s.logger.Error("Admin not found", zap.Error(err))
		return
	}

	s.logger.Info("Admin retrieved", zap.Any("admin", admin))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(admin)
}

func (s *Server) handleUpdateAdmin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		s.logger.Error("Email is required")
		return
	}

	// First, get the existing admin
	existingAdmin, err := s.db.GetAdminByEmail(email)
	if err != nil {
		http.Error(w, "Admin not found", http.StatusNotFound)
		s.logger.Error("Admin not found", zap.Error(err))
		return
	}

	// Only update password if provided in request
	var updateData struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update only the password if provided
	if updateData.Password != "" {
		existingAdmin.Password = updateData.Password
	}

	err = s.db.UpdateAdmin(existingAdmin)
	if err != nil {
		s.logger.Error("Admin update failed", zap.Error(err))
		http.Error(w, "Failed to update admin", http.StatusInternalServerError)
		return
	}

	s.logger.Info("Admin updated", zap.Any("admin", existingAdmin))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(existingAdmin); err != nil {
		s.logger.Error("Failed to encode response", zap.Error(err))
	}
}

func (s *Server) handleDeleteAdmin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		s.logger.Error("Email is required")
		return
	}

	if err := s.db.DeleteAdmin(email); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		s.logger.Error("Admin deletion failed", zap.Error(err))
		return
	}

	s.logger.Info("Admin deleted", zap.Any("email", email))

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := s.signIn(loginRequest.Email, loginRequest.Password)
	if err != nil {
		s.logger.Error("Login failed", zap.Error(err))
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) signIn(email string, password string) (string, error) {

	var err error

	admin, err := s.db.GetAdmin(email)
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(admin.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(admin.ID))
}

func validateId(id string) (int, bool) {
	var intId int
	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, false
	}
	return intId, true
}
