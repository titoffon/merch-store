package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/titoffon/merch-store/internal/db"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct{
	Error string `json:"error"`
}

func Auth(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
			slog.Warn("Username and password must not be empty.")
            return
        }

		if req.Username == "" || req.Password == "" {
			slog.Warn("Validation failed")
			//http.Error(w, "Validation failed", http.StatusBadRequest)
			res, err := json.Marshal(ErrorResponse{
				Error: "Validation failed",
			})
			if err != nil{
				slog.Error("failed Unmarshall")
				http.Error(w, "Internal error", http.StatusInternalServerError  )
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest )
			json.NewEncoder(w).Encode(res)
			return
		}
		
		user, err := conDB.GetUserByName(r.Context(), req.Username)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			slog.Error("Database error", slog.String("error", err.Error()))
			return
		}

		if user == nil {
			err = conDB.CreateUser(r.Context(), req.Username, req.Password)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				slog.Error("Failed to create user", slog.String("error", err.Error()))
				return
			}
			user = &db.User{
				Username: req.Username,
				Balance:  1000,
			}
		} else {
			valid, err := conDB.CheckPassword(r.Context(), req.Username, req.Password)
			if err != nil || !valid {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				slog.Warn("Invalid password")
				return
			}
		}

		token, err := generateJWTToken(user.Username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			slog.Error("Failed to generate token", slog.String("error", err.Error()))
			return
		}

		resp := AuthResponse{
			Token: token,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func generateJWTToken(username string) (string, error) {

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // токен истекает через 72 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GetUserInfo(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserPurchases(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserTransactions(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func PurchaseMerch(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func TransferCoins(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}