package domain

import "time"

type Attendance struct {
	ID           string
	EmployeeID   string
	EmployeeName string
	Location     string
	CheckIn      time.Time
	CheckOut     *time.Time
	IsLate       bool
	Date         time.Time
}

const (
	OfficeStartHour   = 9 // 9 AM
	OfficeStartMinute = 0
)

type CheckInParams struct {
	ID           string
	EmployeeID   string
	EmployeeName string
	Location     string
	CheckInTime  time.Time
}

func NewAttendance(params CheckInParams) *Attendance {
	isLate := false

	limit := time.Date(
		params.CheckInTime.Year(), params.CheckInTime.Month(), params.CheckInTime.Day(),
		OfficeStartHour, OfficeStartMinute, 0, 0, params.CheckInTime.Location(),
	)

	if params.CheckInTime.After(limit) {
		isLate = true
	}

	// Extract date only (year, month, day) for the Date field
	dateOnly := time.Date(
		params.CheckInTime.Year(), params.CheckInTime.Month(), params.CheckInTime.Day(),
		0, 0, 0, 0, params.CheckInTime.Location(),
	)

	return &Attendance{
		ID:           params.ID,
		EmployeeID:   params.EmployeeID,
		EmployeeName: params.EmployeeName,
		Location:     params.Location,
		CheckIn:      params.CheckInTime,
		IsLate:       isLate,
		Date:         dateOnly,
	}
}

func (a *Attendance) SetCheckOut(checkOut time.Time) {
	a.CheckOut = &checkOut
}
