-- Add wallet_address and encrypted_private_key to STUDENT table
ALTER TABLE STUDENT ADD COLUMN IF NOT EXISTS wallet_address VARCHAR(42) UNIQUE;
ALTER TABLE STUDENT ADD COLUMN IF NOT EXISTS encrypted_private_key TEXT;
