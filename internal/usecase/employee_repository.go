package usecase

import (
	"context"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type EmployeeRepository interface {
	Save(ctx context.Context, employee *domain.Employee) error
	FindByID(ctx context.Context, id string) (*domain.Employee, error)
	FindByEmail(ctx context.Context, email string) (*domain.Employee, error)
	FindAll(ctx context.Context) ([]*domain.Employee, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id string) error
}
