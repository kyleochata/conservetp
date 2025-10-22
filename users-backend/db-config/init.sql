-- Create users table directly (no database switching needed)
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

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Optional: Insert sample data
INSERT INTO users (email, name, pwd) 
VALUES 
    ('admin@example.com', 'Admin User', 'hashed_password_123'),
    ('user@example.com', 'John Doe', 'hashed_password_456')
ON CONFLICT (email) DO NOTHING;