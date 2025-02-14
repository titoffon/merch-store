package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/titoffon/merch-store/internal/db"
)

type UserClaims struct {
	Username string `json:"sub"`
	jwt.RegisteredClaims
}

func PurchaseMerch(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		item := chi.URLParam(r, "item")
		if item == "" {
			slog.Error("Item name is required")
			ResponseError(w, http.StatusBadRequest, "Item name is required")
			return
		}

		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			slog.Error("Authorization token is required")
			ResponseError(w, http.StatusUnauthorized, "Authorization token is required")
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := validateJWT(tokenStr, []byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			ResponseError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		username := claims.Username
		user, err := conDB.GetUserByName(r.Context(), username)
		if err != nil || user == nil {
			slog.Error("User by token not found")
			ResponseError(w, http.StatusUnauthorized, "User not found")
			return
		}


		price, err := conDB.GetItemPrice(r.Context(), item)
		if err != nil{
			slog.Error("error with the price", slog.String("error", err.Error()))
			ResponseError(w, 400, "error with the item")
		}

		if user.Balance < price {
			slog.Info("Not enough coins")
			ResponseError(w, http.StatusBadRequest, "Not enough coins")
			return
		}

		tx, err := conDB.DBPool.Begin(r.Context())
		if err != nil {
			slog.Error("Transaction start error")
			ResponseError(w, http.StatusInternalServerError, "Transaction start error")
			return
		}
		defer tx.Rollback(r.Context())

		user.Balance = user.Balance - price
		user, err = conDB.UpdateUserBalance(r.Context(), user, tx)
		if err != nil {
			slog.Error("Transaction start error", slog.String("error", err.Error()))
			ResponseError(w, http.StatusInternalServerError, "Failed to record purchase")
			return
		}
		fmt.Println(user)

		purchase, err := conDB.InsertPurchases(r.Context(), db.Purchases{
				Username: user.Username,
   				Merch_item: item,
			}, tx)
		if err != nil {
			slog.Error("Transaction start error", slog.String("error", err.Error()))
			ResponseError(w, http.StatusInternalServerError, "Failed to record purchase")
			return
		}

		fmt.Println(purchase)

		transactionLog, err := conDB.InsertTransaction_log(r.Context(), db.TransactionLog{
			Sender: user.Username,
			Recipient: "",
			Amount: price,
			TypeOfTransaction: "purchase",
		}, tx)
		if err != nil {
			slog.Error("Transaction start error", slog.String("error", err.Error()))
			ResponseError(w, http.StatusInternalServerError, "Failed to log transaction")
			return
		}

		fmt.Println(transactionLog)

		err = tx.Commit(r.Context())
		if err != nil {
			ResponseError(w, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func validateJWT(tokenStr string, secretKey []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}