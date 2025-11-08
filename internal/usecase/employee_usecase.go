package usecase

import (
	"errors"
	"strings"
	"time"

	domain "github.com/zuyatna/shop-retail-employee-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const maxPhotoSize = 5 * 1024 * 1024 // 5MB

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

func trim(s string) string { return strings.TrimSpace(s) }

func validateRequireContract(employee *domain.Employee) error {
	employee.Email = trim(strings.ToLower(employee.Email))
	employee.PasswordHash = trim(employee.PasswordHash)
	employee.Phone = trim(employee.Phone)

	if employee.Name == "" || employee.Email == "" || employee.Address == "" || employee.District == "" || employee.City == "" || employee.Province == "" || employee.Phone == "" {
		return errors.New("invalid data input, please check data requirement")
	}
	return nil
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

	if err := validateRequireContract(employee); err != nil {
		return err
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

func (u *EmployeeUsecase) Update(callerRole domain.Role, callerID string, employee *domain.Employee) error {
	if employee.ID == "" {
		return domain.ErrBadRequest
	}

	if !(callerRole == domain.RoleSupervisor || callerRole == domain.RoleManager || callerRole == domain.RoleHR) {
		if callerRole == domain.RoleStaff && callerID == employee.ID {
			// Allow staff to update their own employee record
		} else {
			return domain.ErrForbidden
		}
	}

	// validate required fields
	if err := validateRequireContract(employee); err != nil {
		return err
	}

	var existing *domain.Employee
	var err error
	if !employee.PhotoProvided || len(employee.Photo) == 0 {
		existing, err = u.repo.FindByID(employee.ID)
		if err != nil {
			return err
		}
	}

	if !employee.PhotoProvided {
		if existing == nil {
			existing, err = u.repo.FindByID(employee.ID)
			if err != nil {
				return err
			}
		}
		employee.Photo = existing.Photo
	}

	if employee.PhotoProvided && employee.Photo != nil && len(employee.Photo) > maxPhotoSize {
		return domain.ErrPhotoTooLarge
	}

	employee.Email = trim(strings.ToLower(employee.Email))
	if employee.PasswordHash != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		employee.PasswordHash = string(passwordHash)
	} else {
		if existing == nil {
			existing, err = u.repo.FindByID(employee.ID)
			if err != nil {
				return err
			}
		}
		employee.PasswordHash = existing.PasswordHash
	}

	employee.UpdatedAt = time.Now()
	return u.repo.Update(employee)
}

func (u *EmployeeUsecase) UpdatePhoto(callerRole domain.Role, callerID string, id string, photo []byte, photoMIME string) error {
	if !(callerRole == domain.RoleSupervisor || callerRole == domain.RoleManager || callerRole == domain.RoleHR) {
		if !(callerRole == domain.RoleStaff && callerID == id) {
			return domain.ErrForbidden
		}
	}

	if len(photo) == 0 {
		return domain.ErrBadRequest
	}

	if len(photo) > maxPhotoSize {
		return domain.ErrPhotoTooLarge
	}

	switch photoMIME {
	case "image/jpeg", "image/png", "image/gif":
	default:
		return domain.ErrBadRequest
	}

	employee, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	employee.Photo = photo

	if photoMIME == "image/jpg" {
		photoMIME = "image/jpeg"
	}

	employee.PhotoMIME = strings.ToLower(photoMIME)
	employee.UpdatedAt = time.Now()
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
