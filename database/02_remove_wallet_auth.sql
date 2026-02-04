-- Remove wallet_address and encrypted_private_key from STUDENT table if they exist
ALTER TABLE STUDENT DROP COLUMN IF EXISTS wallet_address;
ALTER TABLE STUDENT DROP COLUMN IF EXISTS encrypted_private_key;
