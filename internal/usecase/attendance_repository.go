package usecase

import (
	"context"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type AttendanceRepository interface {
	Save(ctx context.Context, attendance *domain.Attendance) error
	Update(ctx context.Context, attendance *domain.Attendance) error
	FindByEmployeeIDAndDate(ctx context.Context, employeeID string, date time.Time) (*domain.Attendance, error)
}
