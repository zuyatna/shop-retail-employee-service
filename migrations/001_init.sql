-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Employees table
CREATE TABLE employees (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    position VARCHAR(100),
    salary NUMERIC(15,2),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    birthdate DATE,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    phone_number VARCHAR(20) UNIQUE,
    photo TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Status constraint
ALTER TABLE employees
ADD CONSTRAINT chk_employee_status
CHECK (status IN ('active', 'inactive', 'suspended'));

-- Indexes
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_status ON employees(status);
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at);
