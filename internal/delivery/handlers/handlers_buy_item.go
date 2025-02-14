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

func (h *Handlers) PurchaseMerch(w http.ResponseWriter, r *http.Request) {

		item := chi.URLParam(r, "item")
		if item == "" {
			slog.Error("Item name is required")
			ResponseError(w, http.StatusBadRequest, "Item name is required")
			return
		}

		username, err := ExtractJWT(w, r)
		if err != nil {
			return
		}


		price, err := h.Dal.GetItemPrice(r.Context(), item)
		if err != nil{
			slog.Error("error with the price", slog.String("error", err.Error()))
			ResponseError(w, 400, "error with the item")
		}

		tx, err := h.Dal.DBPool.Begin(r.Context())
		if err != nil {
			slog.Error("Transaction start error")
			ResponseError(w, http.StatusInternalServerError, "Transaction start error")
			return
		}
		defer func(){
			txErr := tx.Rollback(r.Context())
			if txErr != nil{
				slog.Info(txErr.Error())
			}
			}()

		err = h.Dal.MinusUserBalance(r.Context(), username, price, tx)
		if err != nil {
			if errors.Is(err, db.ErrLowBalance){
				ResponseError(w, 400, "No enough coins")
				return
			}
			slog.Error("Transaction failed", slog.String("error", err.Error()))
			ResponseError(w, http.StatusInternalServerError, "Transaction failed") 
			return
		}

		err = h.Dal.InsertPurchases(r.Context(), db.Purchases{
				Username: username,
   				Merch_item: item,
			}, tx)
		if err != nil {
			slog.Error("Transaction start error", slog.String("error", err.Error()))
			ResponseError(w, http.StatusInternalServerError, "Failed to record purchase")
			return
		}
/*
		transactionLog, err := h.Dal.InsertTransaction_log(r.Context(), db.TransactionLog{
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
*/
//TO DO изменить БД 

		err = tx.Commit(r.Context())
		if err != nil {
			ResponseError(w, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}


func ExtractJWT(w http.ResponseWriter, r *http.Request) (string, error){

	tokenStr := r.Header.Get("Authorization")

	if tokenStr == "" {
			slog.Error("Authorization token is required")
			ResponseError(w, http.StatusUnauthorized, "Authorization token is required")
			return "", fmt.Errorf("Authorization token is required")
		}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	claims, err := validateJWT(tokenStr, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ResponseError(w, http.StatusUnauthorized, "Invalid token")
		return "", fmt.Errorf("Invalid token")
	}
	if claims.Username == ""{
		ResponseError(w, http.StatusUnauthorized, "Empty Username Plaload")
		slog.Error("Empty Username Plaload")
		return "", fmt.Errorf("Empty Username Plaload")
	}
	return claims.Username, nil
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