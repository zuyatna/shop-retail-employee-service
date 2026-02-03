package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/config"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/attendance"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/clock"
)

type AttendanceUsecase struct {
	attendanceRepo AttendanceRepository
	employeeRepo   EmployeeRepository
	idGen          IDGenerator
	cfg            *config.Config
	clock          clock.Clock
	ctxTimeout     time.Duration
}

func NewAttendanceUsecase(attendanceRepo AttendanceRepository, employeeRepo EmployeeRepository, idGen IDGenerator, cfg *config.Config, clk clock.Clock, timeout time.Duration) *AttendanceUsecase {
	return &AttendanceUsecase{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		idGen:          idGen,
		cfg:            cfg,
		clock:          clk,
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

	now := uc.clock.Now().In(uc.cfg.AppTimezone)
	dateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

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
		ID:              attendanceID,
		EmployeeID:      employeeID,
		EmployeeName:    employee.Name(),
		Location:        req.Location,
		CheckInTime:     now,
		OfficeStartHour: uc.cfg.OfficeStartHour,
		OfficeStartMin:  uc.cfg.OfficeStartMin,
	})

	if err := uc.attendanceRepo.Save(ctx, newAttendance); err != nil {
		return "", fmt.Errorf("failed to save attendance: %w", err)
	}
	slog.Log(ctx, slog.LevelInfo, "Employee checked in", "employeeID", employeeID, "time", now)

	return attendanceID, nil
}

func (uc *AttendanceUsecase) CheckOut(ctx context.Context, employeeID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	now := uc.clock.Now().In(uc.cfg.AppTimezone)
	dateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

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

	attendanceRecord.SetCheckOut(now)

	if err := uc.attendanceRepo.Update(ctx, attendanceRecord); err != nil {
		return fmt.Errorf("failed to update attendance record: %w", err)
	}
	slog.Log(ctx, slog.LevelInfo, "Employee checked out", "employeeID", employeeID, "time", now)

	return nil
}
