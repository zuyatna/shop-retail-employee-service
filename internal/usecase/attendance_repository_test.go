package usecase_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type MockAttendanceRepo struct {
	mock.Mock
}

type MockClock struct {
	currentTime time.Time
}

func (m *MockAttendanceRepo) Save(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) Update(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) FindByEmployeeIDAndDate(ctx context.Context, employeeID string, date time.Time) (*domain.Attendance, error) {
	args := m.Called(ctx, employeeID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attendance), args.Error(1)
}

func (m MockClock) Now() time.Time {
	return m.currentTime
}
