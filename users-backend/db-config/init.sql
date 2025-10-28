-- DROP TABLE IF EXISTS addresses CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

-- Create users table 
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    pwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
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

CREATE OR  REPLACE FUNCTION update_updated_at_col()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_col();

CREATE TRIGGER update_addr_updated_at
    BEFORE UPDATE ON addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_col();

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Fix the sample data insertion
INSERT INTO users (id, email, name, pwd, created_at, updated_at, last_login, is_active) 
VALUES 
    ('adb7f047-6e24-4115-a390-f87c8bd5ab51', 'admin@example.com', 'Admin User', 'hashed_password_123', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, true),
    (gen_random_uuid(), 'user@example.com', 'John Doe', 'hashed_password_456', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, true)
ON CONFLICT (email) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP;