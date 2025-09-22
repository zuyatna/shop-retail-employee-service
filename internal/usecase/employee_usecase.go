package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IDGenerator interface {
	NewID() (string, error)
}

type EmployeeUsecase struct {
	repo  domain.EmployeeRepository
	idGen IDGenerator
}

func NewEmployeeUsecase(repo domain.EmployeeRepository, idGen IDGenerator) *EmployeeUsecase {
	return &EmployeeUsecase{repo: repo, idGen: idGen}
}

func canManageAll(role domain.Role) bool {
	switch role {
	case domain.RoleSupervisor, domain.RoleHR, domain.RoleManager:
		return true
	default:
		return false
	}
}

func (u *EmployeeUsecase) Create(callerRole domain.Role, employee *domain.Employee) error {
	if !canManageAll(callerRole) {
		return domain.ErrForbidden
	}

	employee.Name = strings.TrimSpace(employee.Name)
	employee.Email = strings.TrimSpace(strings.ToLower(employee.Email))
	if employee.Name == "" || employee.Email == "" || employee.PasswordHash == "" {
		return domain.ErrBadRequest
	}

	if employee.ID == "" {
		id, err := u.idGen.NewID()
		if err != nil {
			return err
		}
		employee.ID = id
	}

	_, err := u.repo.FindByEmail(employee.Email)
	if err == nil {
		return domain.ErrDuplicate
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	employee.PasswordHash = string(passwordHash)

	now := time.Now()
	employee.CreatedAt, employee.UpdatedAt = now, now

	return u.repo.Create(employee)
}

func (u *EmployeeUsecase) FindAll(callerRole domain.Role) ([]*domain.Employee, error) {
	if !canManageAll(callerRole) {
		return nil, domain.ErrForbidden
	}
	return u.repo.FindAll()
}

func (u *EmployeeUsecase) FindByID(callerRole domain.Role, callerId, id string) (*domain.Employee, error) {
	if callerRole == domain.RoleStaff && callerId != id {
		return nil, domain.ErrForbidden
	}
	return u.repo.FindByID(id)
}

func (u *EmployeeUsecase) FindByEmail(email string) (*domain.Employee, error) {
	return u.repo.FindByEmail(email)
}

func (u *EmployeeUsecase) Update(employee *domain.Employee) error {
	if employee.ID == "" {
		return domain.ErrBadRequest
	}
	current, err := u.repo.FindByID(employee.ID)
	if err != nil {
		return err
	}
	if employee.PasswordHash == "" {
		employee.PasswordHash = current.PasswordHash
	} else {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		employee.PasswordHash = string(passwordHash)
	}
	return u.repo.Update(employee)
}

func (u *EmployeeUsecase) Delete(callerRole domain.Role, id string) error {
	if !canManageAll(callerRole) {
		return domain.ErrForbidden
	}
	if strings.TrimSpace(id) == "" {
		return domain.ErrBadRequest
	}
	return u.repo.Delete(id)
}
