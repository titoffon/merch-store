package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct{
	DBPool  *pgxpool.Pool
}

type User struct{
	Username string
	HashedPassword string
	Balance     int64
}

type Purchases struct {
	Username string
    Merch_item string
}

type TransactionLog struct {
    Sender string
    Recipient string
    Amount int64
    TypeOfTransaction string
}

func New(ctx context.Context, connectionString string) (*DB, error) {

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	slog.Info("The connection to the database is established")
	return &DB{DBPool: pool}, nil
}

func (r *DB) GetUserByName(ctx context.Context, name string) (*User, error){
	
	q := "SELECT username, hashed_password, balance FROM users WHERE username = $1"
	row := r.DBPool.QueryRow(ctx, q, name)

	var user User
	if err := row.Scan(&user.Username, &user.HashedPassword, &user.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, nil
		}
		
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &user, nil
}

func (r *DB) CreateUser(ctx context.Context, user User) (*User, error) {


	q := "INSERT INTO users (username, hashed_password, balance) VALUES ($1, $2, $3)"
	_, err := r.DBPool.Exec(ctx, q, user.Username, string(user.HashedPassword), user.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return &user, nil
}

func (r *DB) GetItemPrice(ctx context.Context, item string) (int64, error) {
	
	q := "SELECT price FROM merch WHERE name = $1"

	row := r.DBPool.QueryRow(ctx, q, item)

	var price int64 
	if err := row.Scan(&price); err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return 0, fmt.Errorf("There is no such product: %w", err)
		}
		
		return 0, fmt.Errorf("failed to query item: %w", err)
	}
	return price, nil
}

var ErrLowBalance = errors.New("No enough coins")

func (r *DB) MinusUserBalance(ctx context.Context, username string, price int64, tx pgx.Tx) (error){

	q := "UPDATE users SET balance = balance - $1 WHERE username = $2"
	var err error
	if tx == nil{
		_, err = r.DBPool.Exec(ctx, q, price, username )		
	} else {
		_, err = tx.Exec(ctx, q, price, username )
	}
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr){
			if pgErr.ConstraintName == "users_balance_non_negative"{
				return ErrLowBalance
			}
		}
		return fmt.Errorf("failed to update user balance: %w", err)
	}
	return  nil
}

func (r *DB) InsertPurchases(ctx context.Context, purchase Purchases, tx pgx.Tx) (error){

	q := "INSERT INTO purchases (username, merch_item) VALUES ($1, $2)"
	var err error
	if tx == nil{
		_, err = r.DBPool.Exec(ctx, q, purchase.Username, purchase.Merch_item)		
	} else {
		_, err = tx.Exec(ctx, q, purchase.Username, purchase.Merch_item)
	}
	if err != nil {
		return fmt.Errorf("failed to INSERT INTO purchases: %w", err)
	}
	return nil
}

func (r *DB) InsertTransaction_log(ctx context.Context, transaction TransactionLog, tx pgx.Tx) (*TransactionLog, error){

	q := "INSERT INTO transaction_log (sender, recipient, amount, type) VALUES ($1, $2, $3, $4)"
	_, err := r.DBPool.Exec(ctx, q, transaction.Sender, transaction.Recipient, transaction.Amount, transaction.TypeOfTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to INSERT INTO transaction_log: %w", err)
	}

	return &transaction, nil
}