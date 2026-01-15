package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zuyatna/shop-retail-employee-service/internal/adapter/repo/record"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
)

type PostgresEmployeeRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresEmployeeRepo(pool *pgxpool.Pool) *PostgresEmployeeRepo {
	return &PostgresEmployeeRepo{
		pool: pool,
	}
}

func (r *PostgresEmployeeRepo) Save(ctx context.Context, employee *domain.Employee) error {
	rec := record.FromDomain(employee)

	query := `
		INSERT INTO employees (
			id, name, email, password, role, position, salary, status,
			birthdate, address, city, province, phone_number,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			NOW(), NOW()
		)
	`

	_, err := r.pool.Exec(ctx, query,
		rec.ID,
		rec.Name,
		rec.Email,
		rec.Password,
		rec.Role,
		rec.Position,
		rec.Salary,
		rec.Status,
		rec.BirthDate,
		rec.Address,
		rec.City,
		rec.Province,
		rec.PhoneNumber,
	)

	return err
}

func (r *PostgresEmployeeRepo) FindByID(ctx context.Context, id domain.EmployeeID) (*domain.Employee, error) {
	query := `
		SELECT id, name, email, password, role, position, salary, status,
		       birthdate, address, city, province, phone_number, photo,
		       created_at, updated_at, deleted_at
		FROM employees
		WHERE id = $1 AND deleted_at IS NULL
	`

	rows, _ := r.pool.Query(ctx, query, string(id))

	// Use pgx.CollectOneRow to fetch a single row and map it to EmployeeRecord
	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[record.EmployeeRecord]) // pgx.RowToStructByNameLax maps columns to struct fields by name
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

	return rec.ToDomain()
}

func (r *PostgresEmployeeRepo) FindByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	query := `
		SELECT id, name, email, password, role, position, salary, status,
		       birthdate, address, city, province, phone_number, photo,
		       created_at, updated_at, deleted_at
		FROM employees
		WHERE email = $1 AND deleted_at IS NULL
	`

	rows, _ := r.pool.Query(ctx, query, email)

	// Use pgx.CollectOneRow to fetch a single row and map it to EmployeeRecord
	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[record.EmployeeRecord]) // pgx.RowToStructByNameLax maps columns to struct fields by name
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to find employee by email: %w", err)
	}

	return rec.ToDomain()
}

func (r *PostgresEmployeeRepo) Update(ctx context.Context, employee *domain.Employee) error {
	rec := record.FromDomain(employee)

	query := `
		UPDATE employees SET
			name = $1,
			role = $2,
			position = $3,
			salary = $4,
			status = $5,
			birthdate = $6,
			address = $7,
			city = $8,
			province = $9,
			phone_number = $10,
			photo = $11,
			updated_at = NOW()
		WHERE id = $12 AND deleted_at IS NULL
	`

	cmdTag, err := r.pool.Exec(ctx, query,
		rec.Name, rec.Role, rec.Position, rec.Salary, rec.Status,
		rec.BirthDate, rec.Address, rec.City, rec.Province, rec.PhoneNumber, rec.Photo,
		rec.ID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("employee not found or deleted")
	}

	return nil
}

func (r *PostgresEmployeeRepo) Delete(ctx context.Context, id domain.EmployeeID) error {
	query := `
		UPDATE employees
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.pool.Exec(ctx, query, string(id))

	return err
}
