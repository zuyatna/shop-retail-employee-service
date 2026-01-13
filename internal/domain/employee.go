package domain

import (
	"errors"
	"strings"
	"time"
)

type EmployeeID string
type Email string
type Role string
type Status string

const (
	RoleAdmin      Role = "admin"
	RoleSupervisor Role = "supervisor"
	RoleStaff      Role = "staff"
)

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
)

type Employee struct {
	id          EmployeeID
	name        string
	email       Email
	role        Role
	position    string
	salary      int64
	status      Status
	birthdate   *time.Time
	address     string
	city        string
	province    string
	phoneNumber string
	photoURL    string
}

func NewEmployee(
	id EmployeeID,
	name string,
	email Email,
	hashedPassword string,
	role Role,
) (*Employee, error) {

	if id == "" {
		return nil, errors.New("employee ID cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("employee name cannot be empty")
	}

	if !isValidEmail(string(email)) {
		return nil, errors.New("invalid email format")
	}

	if hashedPassword == "" {
		return nil, errors.New("password cannot be empty")
	}

	if !isValidRole(role) {
		return nil, errors.New("invalid role")
	}

	employee := &Employee{
		id:     id,
		name:   name,
		email:  email,
		role:   role,
		status: StatusActive,
	}

	return employee, nil
}

func (e *Employee) Activate() {
	e.status = StatusActive
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func isValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleSupervisor, RoleStaff:
		return true
	default:
		return false
	}
}
