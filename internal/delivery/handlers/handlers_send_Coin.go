package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/titoffon/merch-store/internal/db"
)

type SendCoinRequest struct {
	ToUser  string `json:"toUser"`
	Amount  int64  `json:"amount"`
}

func (h *Handlers) SendCoins(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseError(w, http.StatusBadRequest, "Invalid request body")
		slog.Warn("Invalid request body")
		return
	}

	if req.ToUser == "" || req.Amount <= 0 {
		ResponseError(w, http.StatusBadRequest, "Invalid user or amount")
		return
	}

	username, err := ExtractJWT(w, r)
	if err != nil {
		return
	}

	receiver, err := h.Dal.GetUserByName(r.Context(), req.ToUser)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "Database error")
		slog.Error("Failed to retrieve receiver", slog.String("toUser", req.ToUser), slog.String("error", err.Error()))
		return
	}

	if receiver == nil {
		ResponseError(w, http.StatusBadRequest, "Receiver user does not exist")
		return
	}

	tx, err := h.Dal.DBPool.Begin(r.Context())
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "Transaction start error")
		slog.Error("Failed to start transaction", slog.String("error", err.Error()))
		return
	}
	defer tx.Rollback(r.Context())

	err = h.Dal.MinusUserBalance(r.Context(), username, req.Amount, tx)
	if err != nil {
		if errors.Is(err, db.ErrLowBalance){
			ResponseError(w, 400, "No enough coins")
			return
		}
		ResponseError(w, http.StatusInternalServerError, "Transaction failed")
		slog.Error("Failed to subtract balance", slog.String("error", err.Error()))
		return
	}

	err = h.Dal.PlusUserBalance(r.Context(), req.ToUser, req.Amount, tx)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "Transaction failed")
		slog.Error("Failed to add balance", slog.String("error", err.Error()))
		return
	}

	_, err = h.Dal.InsertTransaction_log(r.Context(), db.TransactionLog{
		Sender:    username,
		Recipient: receiver.Username,
		Amount:    req.Amount,
	}, tx)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "Failed to log transaction")
		slog.Error("Failed to log transaction", slog.String("error", err.Error()))
		return
	}
		
	err = tx.Commit(r.Context())
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}