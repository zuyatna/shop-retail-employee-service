package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/employee"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

func TestEmployeeUsecase_Register(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	mockStorageRepo := new(MockStorageRepo)
	mockIDGen := new(MockIDGenerator)

	ctxTimeout := 2 * time.Second
	uc := usecase.NewEmployeeUsecase(mockRepo, mockStorageRepo, mockIDGen, ctxTimeout)

	req := employee.CreateEmployeeRequest{
		Name:        "Test User",
		Email:       "test@example.com",
		Password:    "password123",
		Role:        "staff",
		Position:    "IT",
		Salary:      1000,
		Status:      "active",
		BirthDate:   "1990-01-01",
		Address:     "Jl Test",
		City:        "Jakarta",
		Province:    "DKI Jakarta",
		PhoneNumber: "08123456789",
	}

	t.Run("Success Register", func(t *testing.T) {
		// 1. Mock Email Check (Not Found / Safe to register)
		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, sql.ErrNoRows).Once()

		// 2. Mock ID Generation
		mockIDGen.On("NewID").Return("uuid-123", nil).Once()

		// 3. Mock Save
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Employee")).Return(nil).Once()

		id, err := uc.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "uuid-123", id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Email Exists", func(t *testing.T) {
		// Mock Email Check (Found / Duplicate)
		existingEmp := &domain.Employee{}
		mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingEmp, nil).Once()

		id, err := uc.Register(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "", id)
		assert.Contains(t, err.Error(), "email already exists")

		// Ensure Save was not called
		mockRepo.AssertNotCalled(t, "Save")
	})
}
