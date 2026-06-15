package repository

import (
	"context"
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// RegisterUserWithWallet menggunakan Database Transaction (Tx)
func (r *UserRepository) RegisterUserWithWallet(nama, email, hashedPassword string) (int, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userID int
	queryUser := `INSERT INTO users (nama, email, password) VALUES ($1, $2, $3) RETURNING id;`
	err = tx.QueryRowContext(ctx, queryUser, nama, email, hashedPassword).Scan(&userID)
	if err != nil {
		return 0, err
	}

	queryWallet := `INSERT INTO wallets (user_id, balance) VALUES ($1, 0.00);`
	_, err = tx.ExecContext(ctx, queryWallet, userID)
	if err != nil {
		return 0, errors.New("gagal membuat dompet digital: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByEmail bertugas mencari user berdasarkan email untuk keperluan login
func (r *UserRepository) GetUserByEmail(email string) (int, string, string, error) {
	query := `SELECT id, nama, password FROM users WHERE email = $1;`
	
	var id int
	var nama, hashedPassword string

	err := r.DB.QueryRow(query, email).Scan(&id, &nama, &hashedPassword)
	if err != nil {
		return 0, "", "", err
	}

	return id, nama, hashedPassword, nil
}

// TopUpWallet bertugas memperbarui saldo di database
func (r *UserRepository) TopUpWallet(userID int, amount float64) (float64, error) {
	ctx := context.Background()
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var newBalance float64
	queryUpdateWallet := `
		UPDATE wallets 
		SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP 
		WHERE user_id = $2 
		RETURNING balance;`
	
	err = tx.QueryRowContext(ctx, queryUpdateWallet, amount, userID).Scan(&newBalance)
	if err != nil {
		return 0, errors.New("gagal mengupdate saldo wallet: " + err.Error())
	}

	var walletID int
	err = tx.QueryRowContext(ctx, "SELECT id FROM wallets WHERE user_id = $1", userID).Scan(&walletID)
	if err != nil {
		return 0, err
	}

	queryLogTx := `
		INSERT INTO transactions (receiver_wallet_id, amount, transaction_type) 
		VALUES ($1, $2, 'topup');`
	
	_, err = tx.ExecContext(ctx, queryLogTx, walletID, amount)
	if err != nil {
		return 0, errors.New("gagal mencatat riwayat transaksi topup: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newBalance, nil
}