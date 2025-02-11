package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func Login(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserInfo(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserPurchases(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserTransactions(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func PurchaseMerch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func TransferCoins(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}