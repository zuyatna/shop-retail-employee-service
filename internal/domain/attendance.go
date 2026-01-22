package domain

import "time"

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
	OfficeStartHour   = 9 // 9 AM
	OfficeStartMinute = 0
)

type CheckInParams struct {
	ID           string
	EmployeeID   string
	EmployeeName string
	Location     string
	CheckInTime  string
}

func NewAttendance(params CheckInParams) *Attendance {
	isLate := false

	layout := time.DateTime
	loc := time.FixedZone("Asia/Jakarta", 7*60*60)
	checkInTime, _ := time.ParseInLocation(layout, params.CheckInTime, loc)

	limit := time.Date(
		checkInTime.Year(), checkInTime.Month(), checkInTime.Day(),
		OfficeStartHour, OfficeStartMinute, 0, 0, checkInTime.Location(),
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
		CheckIn:      params.CheckInTime,
		IsLate:       isLate,
		Date:         dateOnly,
	}
}

func (a *Attendance) SetCheckOut(checkOut string) {
	a.CheckOut = &checkOut
}
