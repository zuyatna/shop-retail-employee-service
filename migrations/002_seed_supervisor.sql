INSERT INTO employees (
    id,
    name,
    email,
    password,
    role,
    position,
    salary,
    status,
    created_at,
    updated_at
) VALUES (
    '018f8c3a-9d9a-7b3e-9b2d-1d0b7c0a0001',
    'Supervisor',
    'supervisor@shop-retail.local',
    '$2a$10$W8kJ6zXv1F2Q9Z9E8ZJZs.0F9J2ZkN1K5qXl1j3q5wPp8P4bKZ7uS',
    'supervisor',
    'Store Supervisor',
    8000000,
    'active',
    NOW(),
    NOW()
);
