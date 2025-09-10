CREATE TABLE IF NOT EXISTS shop_retail_employees (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('supervisor', 'hr', 'manager', 'staff')),
    position TEXT NOT NULL,
    salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Index tambahan untuk pencarian cepat
CREATE INDEX IF NOT EXISTS idx_employees_email ON shop_retail_employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_role ON shop_retail_employees(role);
