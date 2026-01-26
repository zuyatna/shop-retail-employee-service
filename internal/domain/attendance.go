package domain

import (
	"time"
)

type Attendance struct {
	ID           string
	EmployeeID   string
	EmployeeName string
	Location     string
	CheckIn      string
	CheckOut     *string
	IsLate       bool
	Date         time.Time
}

const (
	OfficeStartMinute = 0
)

type CheckInParams struct {
	ID              string
	EmployeeID      string
	EmployeeName    string
	Location        string
	CheckInTime     time.Time
	OfficeStartHour int
	OfficeStartMin  int
}

func NewAttendance(params CheckInParams) *Attendance {
	isLate := false

	checkInTime := params.CheckInTime

	limit := time.Date(
		checkInTime.Year(), checkInTime.Month(), checkInTime.Day(),
		params.OfficeStartHour, OfficeStartMinute, 0, 0, checkInTime.Location(),
	)

	if checkInTime.After(limit) {
		isLate = true
	}

	// Extract date only (year, month, day) for the Date field
	dateOnly := time.Date(
		checkInTime.Year(), checkInTime.Month(), checkInTime.Day(),
		0, 0, 0, 0, checkInTime.Location(),
	)

	return &Attendance{
		ID:           params.ID,
		EmployeeID:   params.EmployeeID,
		EmployeeName: params.EmployeeName,
		Location:     params.Location,
		CheckIn:      checkInTime.Format(time.DateTime),
		IsLate:       isLate,
		Date:         dateOnly,
	}
}

func (a *Attendance) SetCheckOut(checkOut time.Time) {
	t := checkOut.Format(time.DateTime)
	a.CheckOut = &t
}
