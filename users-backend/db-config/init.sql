-- Create users table 
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    pwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true
);

-- Create addresses table. Constraint fk user id
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    street VARCHAR(500) NOT NULL,
    apt_num VARCHAR(50),
    zipcode VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'US',
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_address_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create partial unique index for one primary address per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_primary_address_per_user
ON addresses (user_id)
WHERE is_primary = true;



-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Optional: Insert sample data
INSERT INTO users (email, name, pwd) 
VALUES 
    ('admin@example.com', 'Admin User', 'hashed_password_123'),
    ('user@example.com', 'John Doe', 'hashed_password_456')
ON CONFLICT (email) DO NOTHING;