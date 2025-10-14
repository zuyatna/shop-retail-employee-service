package domain

import (
	"fmt"
	"time"
)

type Role string

const (
	RoleSupervisor Role = "supervisor"
	RoleManager    Role = "manager"
	RoleHR         Role = "hr"
	RoleStaff      Role = "staff"
)

type Employee struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Role          Role       `json:"role"`
	Position      string     `json:"position"`
	Salary        float64    `json:"salary"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Address       string     `json:"address"`
	District      string     `json:"district"`
	City          string     `json:"city"`
	Province      string     `json:"province"`
	Phone         string     `json:"phone"`
	Photo         []byte     `json:"photo,omitempty"`
	PhotoProvided bool       `json:"-"`
}

type EmployeeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type EmployeeRepository interface {
	Create(employee *Employee) error
	FindByID(id string) (*Employee, error)
	FindAll() ([]*Employee, error)
	Update(employee *Employee) error
	Delete(id string) error
	FindByEmail(email string) (*Employee, error)
}

var (
	ErrNotFound      = fmt.Errorf("employee not found")
	ErrDuplicate     = fmt.Errorf("employee already exists")
	ErrBadRequest    = fmt.Errorf("bad request")
	ErrDeleted       = fmt.Errorf("employee has been deleted")
	ErrForbidden     = fmt.Errorf("action forbidden")
	ErrPhotoTooLarge = fmt.Errorf("photo size exceeds the limit")
)
