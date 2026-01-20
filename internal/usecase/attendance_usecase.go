package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		return "", errors.New(EmployeeNotFoundError)
	}

	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

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
	log.Printf("Employee %s checked in at %s \n", employeeID, today.In(time.FixedZone("Asia/Jakarta", 76060)).Format(time.RFC3339))

	return attendanceID, nil
}

func (uc *AttendanceUsecase) CheckOut(ctx context.Context, employeeID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

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
	log.Printf("Attendance %s checked out at %s \n", employeeID, dateOnly.In(time.FixedZone("Asia/Jakarta", 76060)).Format(time.RFC3339))

	return nil
}
