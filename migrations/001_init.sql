CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('supervisor', 'hr', 'manager', 'staff')),
    position TEXT NOT NULL,
    salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL,

    address  TEXT,
    district TEXT,
    city     TEXT,
    province TEXT,
    phone    TEXT,

    -- Foto & MIME
    photo BYTEA,
    photo_mime TEXT,

    -- Constraint ukuran foto (max 5 MB)
    CONSTRAINT employees_photo_max_5mb
      CHECK (photo IS NULL OR octet_length(photo) <= 5 * 1024 * 1024),

    -- Whitelist MIME (jika ada photo, mime harus valid)
    CONSTRAINT employees_photo_mime_whitelist
      CHECK (
        photo IS NULL
        OR photo_mime IN ('image/jpeg','image/png')
      )
);

-- Index tambahan
CREATE INDEX IF NOT EXISTS idx_employees_email       ON employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_role        ON employees(role);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at  ON employees(deleted_at);
