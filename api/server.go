package api

import (
	"employees/internal/db"
	"employees/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.Logger
	router     *mux.Router
	listenAddr string
	db         db.Database
}

func NewServer(listenAddr string, logger *zap.Logger, db db.Database) *Server {
	return &Server{
		logger:     logger,
		listenAddr: listenAddr,
		router:     mux.NewRouter(),
		db:         db,
	}
}

func (s *Server) Start() error {
	s.router.HandleFunc("/employee", s.handleCreateEmployee).Methods("POST")
	s.router.HandleFunc("/employee/{id}", s.handleGetEmployee).Methods("GET")
	s.router.HandleFunc("/employee/{id}", s.handleUpdateEmployee).Methods("PUT")
	s.router.HandleFunc("/employee/{id}", s.handleDeleteEmployee).Methods("DELETE")

	return http.ListenAndServe(s.listenAddr, s.router)
}

func (s *Server) handleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		s.logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.db.CreateEmployee(&emp)
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
	vars := mux.Vars(r)
	id := vars["id"]

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
	vars := mux.Vars(r)
	id := vars["id"]

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
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.DeleteEmployee(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		s.logger.Error("Employee deletion failed", zap.Error(err))
		return
	}

	s.logger.Info("Employee deleted", zap.Any("employeeId", id))

	w.WriteHeader(http.StatusNoContent)
}
