package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/employee"
	"golang.org/x/crypto/bcrypt"
)

type IDGenerator interface {
	NewID() (string, error)
}

type EmployeeUsecase struct {
	repo       EmployeeRepository
	idGen      IDGenerator
	ctxTimeout time.Duration
}

func NewEmployeeUsecase(repo EmployeeRepository, idGen IDGenerator, timeout time.Duration) *EmployeeUsecase {
	return &EmployeeUsecase{
		repo:       repo,
		idGen:      idGen,
		ctxTimeout: timeout,
	}
}

func (uc *EmployeeUsecase) Register(ctx context.Context, req employee.CreateEmployeeRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	existing, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return "", fmt.Errorf("email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := uc.idGen.NewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}

	birthDate, err := parseBirthDate(req.BirthDate)
	if err != nil {
		return "", fmt.Errorf("invalid birth date format: %w", err)
	}

	// Create domain employee
	employee, err := domain.NewEmployee(domain.NewEmployeeParams{
		ID:             domain.EmployeeID(id),
		Name:           req.Name,
		Email:          domain.Email(req.Email),
		HashedPassword: string(hashedBytes),
		Role:           domain.Role(req.Role),
		Position:       req.Position,
		Salary:         int64(req.Salary),
		BirthDate:      birthDate,
		Address:        req.Address,
		City:           req.City,
		Province:       req.Province,
		PhoneNumber:    req.PhoneNumber,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create employee domain: %w", err)
	}

	if err := uc.repo.Save(ctx, employee); err != nil {
		return "", fmt.Errorf("failed to save employee: %w", err)
	}

	return id, nil
}

func parseBirthDate(dateStr string) (*time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
