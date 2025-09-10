package usecase

import (
	"errors"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeUsecase struct {
	repo domain.EmployeeRepository
}

func NewEmployeeUsecase(repo domain.EmployeeRepository) *EmployeeUsecase {
	return &EmployeeUsecase{repo: repo}
}

func (u *EmployeeUsecase) Create(employee *domain.Employee) error {
	if employee.Name == "" || employee.Email == "" || employee.PasswordHash == "" {
		return domain.ErrBadRequest
	}

	_, err := u.repo.FindByEmail(employee.Email)
	if err == nil {
		return domain.ErrDuplicate
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	valid := false
	switch employee.Role {
	case domain.RoleSupervisor, domain.RoleHR, domain.RoleManager, domain.RoleStaff:
		valid = true
	}
	if !valid {
		return domain.ErrBadRequest
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	employee.PasswordHash = string(passwordHash)

	return u.repo.Create(employee)
}

func (u *EmployeeUsecase) FindAll() ([]*domain.Employee, error) {
	return u.repo.FindAll()
}

func (u *EmployeeUsecase) FindByID(id string) (*domain.Employee, error) {
	return u.repo.FindByID(id)
}

func (u *EmployeeUsecase) FindByEmail(email string) (*domain.Employee, error) {
	return u.repo.FindByEmail(email)
}

func (u *EmployeeUsecase) Update(employee *domain.Employee) error {
	if employee.ID == "" {
		return domain.ErrBadRequest
	}
	return u.repo.Update(employee)
}

func (u *EmployeeUsecase) Delete(id string) error {
	if id == "" {
		return domain.ErrBadRequest
	}
	return u.repo.Delete(id)
}
