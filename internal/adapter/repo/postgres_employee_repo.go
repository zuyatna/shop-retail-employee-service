package repo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
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

func (r *PostgresEmployeeRepo) Create(employee *domain.Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `INSERT INTO employees (id, name, email, password_hash, role, position, salary, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		employee.ID, employee.Name, employee.Email, employee.PasswordHash,
		employee.Role, employee.Position, employee.Salary, employee.Status,
		time.Now(), time.Now(),
	)
	if err != nil {
		log.Println("Error inserting employee:", err)
		return err
	}

	res := domain.EmployeeResponse{
		ID:    employee.ID,
		Name:  employee.Name,
		Email: employee.Email,
	}

	log.Println("Employee created: ", res)
	return nil
}

func (r *PostgresEmployeeRepo) FindByID(id string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `SELECT id, name, email, password_hash, role, position, salary, status, created_at, updated_at
			FROM employees WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	employee := &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.PasswordHash,
		&employee.Role, &employee.Position, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Println("Error finding employee by ID:", err)
		return nil, err
	}
	return employee, nil
}

func (r *PostgresEmployeeRepo) FindAll() ([]*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `SELECT id, name, email, password_hash, role, position, salary, status, created_at, updated_at
			  FROM employees`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Println("Error finding all employees:", err)
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		employee := &domain.Employee{}
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.PasswordHash,
			&employee.Role, &employee.Position, &employee.Salary, &employee.Status,
			&employee.CreatedAt, &employee.UpdatedAt)
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

func (r *PostgresEmployeeRepo) FindByEmail(email string) (*domain.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `SELECT id, name, email, password_hash, role, position, salary, status, created_at, updated_at
			  FROM employees WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)
	employee := &domain.Employee{}
	err := row.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.PasswordHash,
		&employee.Role, &employee.Position, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Println("Error finding employee by email:", err)
		return nil, err
	}
	return employee, nil
}

func (r *PostgresEmployeeRepo) Update(employee *domain.Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `UPDATE employees
			  SET name = $1, email = $2, password_hash = $3, role = $4, position = $5, salary = $6, status = $7, updated_at = $8
			  WHERE id = $9`
	_, err := r.pool.Exec(ctx, query,
		employee.Name, employee.Email, employee.PasswordHash,
		employee.Role, employee.Position, employee.Salary, employee.Status,
		time.Now(), employee.ID,
	)
	if err != nil {
		log.Println("Error updating employee:", err)
		return err
	}
	return nil
}

func (r *PostgresEmployeeRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `DELETE FROM employees WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error deleting employee:", err)
		return err
	}
	return nil
}
