package usecase_test

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type MockEmployeeRepo struct {
	mock.Mock
}

type MockStorageRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) Save(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepo) FindByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) FindByID(ctx context.Context, id string) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) FindAll(ctx context.Context) ([]*domain.Employee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) Update(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorageRepo) UploadFile(ctx context.Context, fileName string, contentType string, content io.Reader, size int64) (string, error) {
	args := m.Called(ctx, fileName, contentType, content, size)
	return args.String(0), args.Error(1)
}
