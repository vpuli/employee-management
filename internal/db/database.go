package db

import "employees/internal/models"

//go:generate mockery --name Database
type Database interface {
	CreateEmployee(emp *models.Employee) error
	GetEmployee(id string) (*models.Employee, error)
	UpdateEmployee(emp *models.Employee) error
	DeleteEmployee(id string) error
	Close() error
}
