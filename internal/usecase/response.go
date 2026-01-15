package usecase

import (
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type EmployeeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Position    string     `json:"position"`
	Salary      int64      `json:"salary"`
	Status      string     `json:"status"`
	BirthDate   *time.Time `json:"birth_date"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	PhoneNumber string     `json:"phone_number"`
	Photo       string     `json:"photo,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
}

// FromDomain maps domain.Employee to EmployeeResponse
func FromDomain(e *domain.Employee) *EmployeeResponse {
	if e == nil {
		return nil
	}

	return &EmployeeResponse{
		ID:          string(e.ID()),
		Name:        e.Name(),
		Email:       string(e.Email()),
		Role:        string(e.Role()),
		Position:    e.Position(),
		Salary:      e.Salary(),
		Status:      string(e.Status()),
		BirthDate:   e.BirthDate(),
		Address:     e.Address(),
		City:        e.City(),
		Province:    e.Province(),
		PhoneNumber: string(e.PhoneNumber()),
		Photo:       e.Photo(),
	}
}
