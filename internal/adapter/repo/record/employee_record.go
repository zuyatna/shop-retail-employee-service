package record

import (
	"database/sql"
	"time"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type EmployeeRecord struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Email       string         `db:"email"`
	Password    string         `db:"password"`
	Role        string         `db:"role"`
	Position    sql.NullString `db:"position"`
	Salary      sql.NullInt64  `db:"salary"`
	Status      string         `db:"status"`
	BirthDate   sql.NullTime   `db:"birthdate"`
	Address     sql.NullString `db:"address"`
	City        sql.NullString `db:"city"`
	Province    sql.NullString `db:"province"`
	PhoneNumber sql.NullString `db:"phone_number"`
	Photo       sql.NullString `db:"photo"`

	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

// FromDomain converts a domain.Employee to EmployeeRecord.
func FromDomain(e *domain.Employee) *EmployeeRecord {
	return &EmployeeRecord{
		ID:          string(e.ID()),
		Name:        e.Name(),
		Email:       string(e.Email()),
		Password:    e.PasswordHash(),
		Role:        string(e.Role()),
		Position:    toNullString(e.Position()),
		Salary:      toNullInt64(e.Salary()),
		Status:      string(e.Status()),
		BirthDate:   toNullTime(e.BirthDate()),
		Address:     toNullString(e.Address()),
		City:        toNullString(e.City()),
		Province:    toNullString(e.Province()),
		PhoneNumber: toNullString(e.PhoneNumber()),
		Photo:       toNullString(e.Photo()),
	}
}

// ToDomain converts an EmployeeRecord to domain.Employee.
func (r *EmployeeRecord) ToDomain() (*domain.Employee, error) {
	return domain.ReconstituteEmployee(domain.ReconstituteEmployeeParams{
		ID:           r.ID,
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
		Role:         r.Role,
		Position:     r.Position.String,
		Salary:       r.Salary.Int64,
		Status:       r.Status,
		BirthDate:    validTimeOrNil(r.BirthDate),
		Address:      r.Address.String,
		City:         r.City.String,
		Province:     r.Province.String,
		PhoneNumber:  r.PhoneNumber.String,
		Photo:        r.Photo.String,
	})
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: i != 0}
}

func toNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func validTimeOrNil(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
