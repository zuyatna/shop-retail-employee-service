package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/attendance"
)

type AttendanceUsecase struct {
	attendanceRepo AttendanceRepository
	employeeRepo   EmployeeRepository
	idGen          IDGenerator
	ctxTimeout     time.Duration
}

func NewAttendanceUsecase(attendanceRepo AttendanceRepository, employeeRepo EmployeeRepository, idGen IDGenerator, timeout time.Duration) *AttendanceUsecase {
	return &AttendanceUsecase{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		idGen:          idGen,
		ctxTimeout:     timeout,
	}
}

func (uc *AttendanceUsecase) CheckIn(ctx context.Context, employeeID string, req attendance.CheckInRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	employee, err := uc.employeeRepo.FindByID(ctx, employeeID)
	if err != nil {
		return "", fmt.Errorf("failed to find employee by id: %w", err)
	}
	if employee == nil {
		return "", EmployeeNotFoundError
	}

	today, dateOnly := uc.getJakartaTimeAndDate()

	existingAttendance, err := uc.attendanceRepo.FindByEmployeeIDAndDate(ctx, employeeID, dateOnly)
	if err != nil {
		return "", err
	}
	if existingAttendance != nil {
		return "", errors.New("you've already checked in today")
	}

	attendanceID, err := uc.idGen.NewID()
	if err != nil {
		return "", err
	}

	newAttendance := domain.NewAttendance(domain.CheckInParams{
		ID:           attendanceID,
		EmployeeID:   employeeID,
		EmployeeName: employee.Name(),
		Location:     req.Location,
		CheckInTime:  today,
	})

	if err := uc.attendanceRepo.Save(ctx, newAttendance); err != nil {
		return "", fmt.Errorf("failed to save attendance: %w", err)
	}
	slog.Log(ctx, slog.LevelInfo, "Employee checked in", "employeeID", employeeID, "time", today)

	return attendanceID, nil
}

func (uc *AttendanceUsecase) CheckOut(ctx context.Context, employeeID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	today, dateOnly := uc.getJakartaTimeAndDate()

	attendanceRecord, err := uc.attendanceRepo.FindByEmployeeIDAndDate(ctx, employeeID, dateOnly)
	if err != nil {
		return fmt.Errorf("failed to find attendance record: %w", err)
	}
	if attendanceRecord == nil {
		return errors.New("no check-in record found for today")
	}
	if attendanceRecord.CheckOut != nil {
		return errors.New("you have already checked out today")
	}

	attendanceRecord.SetCheckOut(today)

	if err := uc.attendanceRepo.Update(ctx, attendanceRecord); err != nil {
		return fmt.Errorf("failed to update attendance record: %w", err)
	}
	slog.Log(ctx, slog.LevelInfo, "Employee checked out", "employeeID", employeeID, "time", today)

	return nil
}

func (uc *AttendanceUsecase) getJakartaTimeAndDate() (string, time.Time) {
	now := time.Now()
	today := now.In(time.FixedZone("Asia/Jakarta", 7*60*60)).Format(time.DateTime)
	dateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return today, dateOnly
}
