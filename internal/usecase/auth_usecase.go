package usecase

import (
	"errors"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type TokenSigner interface {
	Generate(userID, email, role string) (string, error)
}

type AuthUsecase struct {
	repo   domain.EmployeeRepository
	signer TokenSigner
}

func NewAuthUsecase(repo domain.EmployeeRepository, signer TokenSigner) *AuthUsecase {
	return &AuthUsecase{repo: repo, signer: signer}
}

func (u *AuthUsecase) Login(email, password string) (string, *domain.Employee, error) {
	if email == "" || password == "" {
		return "", nil, domain.ErrBadRequest
	}

	employee, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := u.signer.Generate(employee.ID, employee.Email, string(employee.Role))
	if err != nil {
		return "", nil, err
	}

	return token, employee, nil
}
