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
	id           EmployeeID
	name         string
	email        Email
	passwordHash string
	role         Role
	position     string
	salary       int64
	status       Status
	birthdate    *time.Time
	address      string
	city         string
	province     string
	phoneNumber  string
	photo        string
}

type NewEmployeeParams struct {
	ID             EmployeeID
	Name           string
	Email          Email
	HashedPassword string
	Role           Role
	Position       string
	Salary         int64
	BirthDate      *time.Time
	Address        string
	City           string
	Province       string
	PhoneNumber    string
}

func NewEmployee(params NewEmployeeParams) (*Employee, error) {

	if params.ID == "" {
		return nil, errors.New("employee ID cannot be empty")
	}

	if strings.TrimSpace(params.Name) == "" {
		return nil, errors.New("employee name cannot be empty")
	}

	if !isValidEmail(string(params.Email)) {
		return nil, errors.New("invalid email format")
	}

	if params.HashedPassword == "" {
		return nil, errors.New("password cannot be empty")
	}

	if !isValidRole(params.Role) {
		return nil, errors.New("invalid role")
	}

	if params.Salary < 0 {
		return nil, errors.New("salary cannot be negative")
	}

	employee := &Employee{
		id:           params.ID,
		name:         params.Name,
		email:        params.Email,
		passwordHash: params.HashedPassword,
		role:         params.Role,
		status:       StatusActive,
		position:     params.Position,
		salary:       params.Salary,
		birthdate:    params.BirthDate,
		address:      params.Address,
		city:         params.City,
		province:     params.Province,
		phoneNumber:  params.PhoneNumber,
	}

	return employee, nil
}

func (e *Employee) Activate() {
	e.status = StatusActive
}

func (e *Employee) Suspend() {
	e.status = StatusSuspended
}

func (e *Employee) Deactivate() {
	e.status = StatusInactive
}

func (e *Employee) ChangeRole(newRole Role) error {
	if !isValidRole(newRole) {
		return errors.New("invalid role")
	}
	e.role = newRole
	return nil
}

type UpdateProfileParams struct {
	Position    string
	Salary      int64
	Address     string
	City        string
	Province    string
	PhoneNumber string
	Photo       string
	BirthDate   *time.Time
}

func (e *Employee) UpdateProfile(params UpdateProfileParams) error {
	if params.Salary < 0 {
		return errors.New("salary cannot be negative")
	}

	e.position = params.Position
	e.salary = params.Salary
	e.address = params.Address
	e.city = params.City
	e.province = params.Province
	e.phoneNumber = params.PhoneNumber
	e.photo = params.Photo
	e.birthdate = params.BirthDate

	return nil
}

func (e *Employee) ID() EmployeeID {
	return e.id
}

func (e *Employee) Name() string {
	return e.name
}

func (e *Employee) Email() Email {
	return e.email
}

func (e *Employee) PasswordHash() string {
	return e.passwordHash
}

func (e *Employee) Role() Role {
	return e.role
}

func (e *Employee) Status() Status {
	return e.status
}

func (e *Employee) Salary() int64 {
	return e.salary
}

func (e *Employee) Position() string {
	return e.position
}

func (e *Employee) BirthDate() *time.Time {
	return e.birthdate
}

func (e *Employee) Address() string {
	return e.address
}

func (e *Employee) City() string {
	return e.city
}

func (e *Employee) Province() string {
	return e.province
}

func (e *Employee) PhoneNumber() string {
	return e.phoneNumber
}

func (e *Employee) Photo() string {
	return e.photo
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

type ReconstituteEmployeeParams struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	Position     string
	Salary       int64
	Status       string
	BirthDate    *time.Time
	Address      string
	City         string
	Province     string
	PhoneNumber  string
	Photo        string
}

func ReconstituteEmployee(p ReconstituteEmployeeParams) (*Employee, error) {
	return &Employee{
		id:           EmployeeID(p.ID),
		name:         p.Name,
		email:        Email(p.Email),
		passwordHash: p.PasswordHash,
		role:         Role(p.Role),
		position:     p.Position,
		salary:       p.Salary,
		status:       Status(p.Status),
		birthdate:    p.BirthDate,
		address:      p.Address,
		city:         p.City,
		province:     p.Province,
		phoneNumber:  p.PhoneNumber,
		photo:        p.Photo,
	}, nil
}
