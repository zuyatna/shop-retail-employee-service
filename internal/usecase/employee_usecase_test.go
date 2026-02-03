package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/attendance"
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

func TestAttendanceUsecase_CheckIn(t *testing.T) {
	mockAttRepo := new(MockAttendanceRepo)
	mockEmpRepo := new(MockEmployeeRepo)
	mockIDGen := new(MockIDGenerator)

	// Setup Config & Timezone
	loc, _ := time.LoadLocation("Asia/Jakarta")
	cfg := &config.Config{
		AppTimezone:     loc,
		OfficeStartHour: 9, // 9 AM
		OfficeStartMin:  0,
	}

	simulateDate := time.Date(2026, 10, 10, 0, 0, 0, 0, loc) // June 10, 2026 00:00:00

	req := attendance.CheckInRequest{
		Location: "Office HQ",
	}
	employeeID := "emp-123"

	t.Run("Success - On Time", func(t *testing.T) {
		mockTime := time.Date(2026, 10, 10, 8, 55, 0, 0, loc) // June 10, 2026 08:55:00
		mockClock := MockClock{currentTime: mockTime}

		uc := usecase.NewAttendanceUsecase(mockAttRepo, mockEmpRepo, mockIDGen, cfg, mockClock, time.Second)

		emp := &domain.Employee{}
		mockEmpRepo.On("FindByID", mock.Anything, employeeID).Return(emp, nil).Once()

		mockAttRepo.On("FindByEmployeeIDAndDate", mock.Anything, employeeID, simulateDate).Return(nil, nil).Once()

		mockIDGen.On("NewID").Return("att-123", nil).Once()

		mockAttRepo.On("Save", mock.Anything, mock.MatchedBy(func(a *domain.Attendance) bool {
			return a.IsLate == false && a.EmployeeID == employeeID
		})).Return(nil).Once()

		id, err := uc.CheckIn(context.Background(), employeeID, req)

		assert.NoError(t, err)
		assert.Equal(t, "att-123", id)
		mockEmpRepo.AssertExpectations(t)
		mockAttRepo.AssertExpectations(t)
		mockIDGen.AssertExpectations(t)
	})

	t.Run("Success - Late", func(t *testing.T) {
		mockTime := time.Date(2026, 10, 10, 9, 15, 0, 0, loc) // June 10, 2026 09:15:00
		mockClock := MockClock{currentTime: mockTime}

		uc := usecase.NewAttendanceUsecase(mockAttRepo, mockEmpRepo, mockIDGen, cfg, mockClock, time.Second)

		emp := &domain.Employee{}
		mockEmpRepo.On("FindByID", mock.Anything, employeeID).Return(emp, nil).Once()

		mockAttRepo.On("FindByEmployeeIDAndDate", mock.Anything, employeeID, simulateDate).Return(nil, nil).Once()

		mockIDGen.On("NewID").Return("att-124", nil).Once()

		mockAttRepo.On("Save", mock.Anything, mock.MatchedBy(func(a *domain.Attendance) bool {
			return a.IsLate == true && a.EmployeeID == employeeID
		})).Return(nil).Once()

		id, err := uc.CheckIn(context.Background(), employeeID, req)

		assert.NoError(t, err)
		assert.Equal(t, "att-124", id)
		mockEmpRepo.AssertExpectations(t)
		mockAttRepo.AssertExpectations(t)
		mockIDGen.AssertExpectations(t)
	})

	t.Run("Fail - Already Checked In", func(t *testing.T) {
		mockTime := time.Date(2026, 10, 10, 8, 55, 0, 0, loc) // June 10, 2026 08:55:00
		mockClock := MockClock{currentTime: mockTime}

		uc := usecase.NewAttendanceUsecase(mockAttRepo, mockEmpRepo, mockIDGen, cfg, mockClock, time.Second)

		emp := &domain.Employee{}
		mockEmpRepo.On("FindByID", mock.Anything, employeeID).Return(emp, nil).Once()

		existingAttendance := &domain.Attendance{}
		mockAttRepo.On("FindByEmployeeIDAndDate", mock.Anything, employeeID, simulateDate).Return(existingAttendance, nil).Once()

		id, err := uc.CheckIn(context.Background(), employeeID, req)

		assert.Error(t, err)
		assert.Equal(t, "", id)
		assert.Contains(t, err.Error(), "you've already checked in today")
		mockEmpRepo.AssertExpectations(t)
		mockAttRepo.AssertExpectations(t)
	})
}
