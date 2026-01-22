package record

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type EmployeeRecord struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Email       string         `db:"email"`
	Password    string         `db:"password"`
	Role        string         `db:"role"`
	Position    sql.NullString `db:"position"`
	Salary      pgtype.Numeric `db:"salary"`
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
		Salary:      int64ToNumeric(e.Salary()),
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
	salary, err := numericToInt64(r.Salary)
	if err != nil {
		return nil, fmt.Errorf("failed to convert salary: %w", err)
	}

	return domain.ReconstituteEmployee(domain.ReconstituteEmployeeParams{
		ID:           r.ID,
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
		Role:         r.Role,
		Position:     r.Position.String,
		Salary:       salary,
		Status:       r.Status,
		BirthDate:    validTimeOrNil(r.BirthDate),
		Address:      r.Address.String,
		City:         r.City.String,
		Province:     r.Province.String,
		PhoneNumber:  r.PhoneNumber.String,
		Photo:        r.Photo.String,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	})
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func int64ToNumeric(v int64) pgtype.Numeric {
	n := pgtype.Numeric{}
	// Represent as an integer numeric (scale 0\)
	n.Int = big.NewInt(v)
	n.Exp = 0
	n.Valid = true
	return n
}

func numericToInt64(n pgtype.Numeric) (int64, error) {
	if !n.Valid {
		return 0, nil
	}
	if n.Int == nil {
		return 0, fmt.Errorf("numeric has nil Int")
	}
	// If salary stored with decimals, this will truncate by shifting according to Exp
	x := new(big.Int).Set(n.Int)
	if n.Exp < 0 {
		den := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
		x.Quo(x, den)
	} else if n.Exp > 0 {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)
		x.Mul(x, mul)
	}
	if !x.IsInt64() {
		return 0, fmt.Errorf("numeric out of int64 range")
	}
	return x.Int64(), nil
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
