package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo       EmployeeRepository
	jwtSigner  *jwtutil.Signer
	ctxTimeout time.Duration
}

func NewAuthUsecase(repo EmployeeRepository, jwtSigner *jwtutil.Signer, timeout time.Duration) *AuthUsecase {
	return &AuthUsecase{
		repo:       repo,
		jwtSigner:  jwtSigner,
		ctxTimeout: timeout,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("login failed: user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	if user.Status() != domain.StatusActive {
		return "", fmt.Errorf("user account is not active")
	}

	token, err := uc.jwtSigner.Generate(string(user.ID()), string(user.Email()), string(user.Role()))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
