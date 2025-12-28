-- Add Google authentication columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS picture_url TEXT,
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
