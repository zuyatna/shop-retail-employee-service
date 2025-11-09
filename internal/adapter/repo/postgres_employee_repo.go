package repo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/zuyatna/shop-retail-employee-service/internal/model"
)

type PostgresEmployeeRepo struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func NewPostgresEmployeeRepo(pool *pgxpool.Pool, timeout time.Duration) *PostgresEmployeeRepo {
	return &PostgresEmployeeRepo{
		pool:    pool,
		timeout: timeout,
	}
}

func (r *PostgresEmployeeRepo) Create(ct context.Context, employee *domain.Employee) error {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `INSERT INTO employees (id, name, email, password_hash, role, position, salary, status, created_at, updated_at, deleted_at, address, district, city, province, phone, photo, photo_mime)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULL, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.pool.Exec(ctx, query,
		employee.ID, employee.Name, employee.Email, employee.PasswordHash,
		employee.Role, employee.Position, employee.Salary, employee.Status,
		employee.CreatedAt, employee.UpdatedAt, employee.Address, employee.District,
		employee.City, employee.Province, employee.Phone, employee.Photo, employee.PhotoMIME,
	)
	if err != nil {
		log.Println("Error inserting employee:", err)
		return err
	}

	res := domain.EmployeeResponse{
		ID:    employee.ID,
		Name:  employee.Name,
		Email: employee.Email,
		Role:  employee.Role,
	}
	log.Println("Employee created: ", res)
	return nil
}

func (r *PostgresEmployeeRepo) FindByID(ct context.Context, id string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `SELECT id, name, email, password_hash, role, position, salary, status, created_at, updated_at, deleted_at,
                    COALESCE(address, ''), COALESCE(district, ''), COALESCE(city, ''), COALESCE(province, ''), COALESCE(phone, ''),
                    COALESCE(photo, ''::bytea), COALESCE(photo_mime, '')
            FROM employees WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	employee := &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.PasswordHash,
		&employee.Role, &employee.Position, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt, &employee.DeletedAt, &employee.Address,
		&employee.District, &employee.City, &employee.Province, &employee.Phone, &employee.Photo, &employee.PhotoMIME)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Println("Error finding employee by ID:", err)
		return nil, err
	}

	if employee.DeletedAt != nil {
		return nil, domain.ErrDeleted
	}
	return employee, nil
}

func (r *PostgresEmployeeRepo) FindAll(ct context.Context) ([]*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `SELECT id, name, email, role, position, salary, status, created_at, updated_at, deleted_at,
                    COALESCE(address, ''), COALESCE(district, ''), COALESCE(city, ''), COALESCE(province, ''), COALESCE(phone, '')
              FROM employees
              WHERE deleted_at IS NULL
              ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Println("Error finding all employees:", err)
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		employee := &domain.Employee{}
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Role,
			&employee.Position, &employee.Salary, &employee.Status, &employee.CreatedAt,
			&employee.UpdatedAt, &employee.DeletedAt, &employee.Address, &employee.District,
			&employee.City, &employee.Province, &employee.Phone)
		if err != nil {
			log.Println("Error scanning employee:", err)
			return nil, err
		}
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over employees:", err)
		return nil, err
	}
	return employees, nil
}

func (r *PostgresEmployeeRepo) FindByEmail(ct context.Context, email string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `SELECT id, name, email, password_hash, role, position, salary, status, created_at, updated_at, deleted_at,
                    COALESCE(address, ''), COALESCE(district, ''), COALESCE(city, ''), COALESCE(province, ''), COALESCE(phone, ''),
                    COALESCE(photo, ''::bytea), COALESCE(photo_mime, '')
              FROM employees WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, email)
	employee := &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.PasswordHash,
		&employee.Role, &employee.Position, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt, &employee.DeletedAt, &employee.Address,
		&employee.District, &employee.City, &employee.Province, &employee.Phone, &employee.Photo, &employee.PhotoMIME)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Println("Error finding employee by email:", err)
		return nil, err
	}

	if employee.DeletedAt != nil {
		return nil, domain.ErrDeleted
	}
	return employee, nil
}

func (r *PostgresEmployeeRepo) Update(ct context.Context, employee *domain.Employee) error {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `UPDATE employees
			  SET name = $1, email = $2, password_hash = $3, role = $4, position = $5, salary = $6, status = $7, updated_at = $8, address = $9, district = $10, city = $11, province = $12, phone = $13, photo = $14, photo_mime = $15
			  WHERE id = $16 AND deleted_at IS NULL`
	cmd, err := r.pool.Exec(ctx, query,
		employee.Name, employee.Email, employee.PasswordHash, employee.Role,
		employee.Position, employee.Salary, employee.Status, time.Now(),
		employee.Address, employee.District, employee.City, employee.Province,
		employee.Phone, employee.Photo, employee.PhotoMIME, employee.ID,
	)
	if err != nil {
		log.Println("Error updating employee:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		_, err := r.FindByID(ctx, employee.ID)
		if err == domain.ErrDeleted {
			return domain.ErrDeleted
		}
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresEmployeeRepo) Delete(ct context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ct, r.timeout)
	defer cancel()

	query := `UPDATE employees
			  SET deleted_at = $1, updated_at = $2
			  WHERE id = $3 AND deleted_at IS NULL`
	cmd, err := r.pool.Exec(ctx, query, time.Now(), time.Now(), id)
	if err != nil {
		log.Println("Error deleting employee:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		_, err := r.FindByID(ctx, id)
		if err == domain.ErrDeleted {
			return domain.ErrDeleted
		}
		return domain.ErrNotFound
	}
	return nil
}
