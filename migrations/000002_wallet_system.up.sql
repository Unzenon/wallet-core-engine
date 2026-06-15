-- 1. Tambah kolom password ke tabel users yang sudah ada
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255) NOT NULL DEFAULT '';

-- 2. Buat tabel wallets (Dompet Digital)
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Buat tabel transactions (Buku Kas/Ledger)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INT,                  -- NULL kalau tipenya Top Up
    receiver_wallet_id INT,                -- NULL kalau tipenya Withdraw/Tarik tunai
    amount NUMERIC(15, 2) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL, -- 'topup', 'transfer'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sender FOREIGN KEY(sender_wallet_id) REFERENCES wallets(id),
    CONSTRAINT fk_receiver FOREIGN KEY(receiver_wallet_id) REFERENCES wallets(id)
);