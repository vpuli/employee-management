package api

import (
	"employees/internal/db"
	"employees/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
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
	s.router.HandleFunc("/employee", s.handleEmployee)
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

func (s *Server) handleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		s.logger.Error("Employee creation failed", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	maxID, err := s.db.GetMaxEmployeeID()
	if err != nil {
		s.logger.Error("Failed to get max employee ID", zap.Error(err))
		http.Error(w, "Employee creation failed", http.StatusInternalServerError)
		return
	}

	nextId, err := strconv.Atoi(maxID)
	if err != nil {
		s.logger.Error("Failed to convert max employee ID to int", zap.Error(err))
		http.Error(w, "Employee creation failed", http.StatusInternalServerError)
		return
	}

	nextId += 1
	emp.ID = strconv.Itoa(nextId)

	err = s.db.CreateEmployee(&emp)
	if err != nil {
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

	if !validateId(id) {
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

	emp.ID = id // Ensure ID matches URL parameter
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

func validateId(id string) bool {
	if _, err := strconv.Atoi(id); err != nil {
		return false
	}
	return true
}
