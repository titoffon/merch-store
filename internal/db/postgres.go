package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DB struct{
	DBPool  *pgxpool.Pool
}

type User struct{
	Username string
	HashedPassword string
	Balance     int64
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
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &user, nil
}

func (r *DB) CreateUser(ctx context.Context, name, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	q := "INSERT INTO users (username, hashed_password, balance) VALUES ($1, $2, $3)"
	_, err = r.DBPool.Exec(ctx, q, name, string(hashedPassword), 1000)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *DB) CheckPassword(ctx context.Context, name, password string) (bool, error) {
	user, err := r.GetUserByName(ctx, name)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		slog.Info("The password is uncorrect")
		return false, nil
	}
	slog.Info("The password is correct")
	return true, nil
}



