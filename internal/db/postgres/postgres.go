package postgres

import (
	"employees/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Employee{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := db.AutoMigrate(&models.Admin{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) CreateEmployee(emp *models.Employee) error {
	return p.db.Create(emp).Error
}

func (p *PostgresDB) GetEmployee(id string) (*models.Employee, error) {
	var emp models.Employee
	if err := p.db.First(&emp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (p *PostgresDB) UpdateEmployee(emp *models.Employee) error {
	return p.db.Save(emp).Error
}

func (p *PostgresDB) DeleteEmployee(id string) error {
	return p.db.Delete(&models.Employee{}, "id = ?", id).Error
}

func (p *PostgresDB) CreateAdmin(admin *models.Admin) error {
	return p.db.Create(admin).Error
}

func (p *PostgresDB) GetAdmin(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := p.db.First(&admin, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (p *PostgresDB) UpdateAdmin(admin *models.Admin) error {
	return p.db.Save(admin).Error
}

func (p *PostgresDB) DeleteAdmin(email string) error {
	return p.db.Delete(&models.Admin{}, "email = ?", email).Error
}

func (p *PostgresDB) GetAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := p.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
