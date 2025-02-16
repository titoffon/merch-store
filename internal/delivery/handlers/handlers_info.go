package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type InfoResponse struct {
    Coins      int64       `json:"coins"`
    Inventory  []InvItem   `json:"inventory"`
    CoinHistory CoinHistory `json:"coinHistory"`
}

type InvItem struct {
    Type     string `json:"type"`
    Quantity int64  `json:"quantity"`
}

type CoinHistory struct {
    Received []ReceivedTx `json:"received"`
    Sent     []SentTx     `json:"sent"`
}

type ReceivedTx struct {
    FromUser string `json:"fromUser"`
    Amount   int64  `json:"amount"`
}

type SentTx struct {
    ToUser string `json:"toUser"`
    Amount int64  `json:"amount"`
}

func (h *Handlers) UserInfo(w http.ResponseWriter, r *http.Request) {
	username, err := ExtractJWT(w, r)
		if err != nil {
		return
	}

	user, err := h.Dal.GetUserByName(r.Context(), username)
    if err != nil {
        slog.Error("Failed to get user by name", slog.String("error", err.Error()))
        ResponseError(w, http.StatusInternalServerError, "Failed to get user by name")
        return
    }

	purchases, err := h.Dal.GetUserPurchases(r.Context(), username)
    if err != nil {
        slog.Error("Failed to get user purchases", slog.String("error", err.Error()))
        ResponseError(w, http.StatusInternalServerError, "Failed to get user purchases")
        return
    }

    var inventory []InvItem
    for _, p := range purchases {
        inventory = append(inventory, InvItem{
            Type:     p.MerchItem,
            Quantity: p.Quantity,
        })
    }

	receivedTxs, err := h.Dal.GetTransactionsReceived(r.Context(), username)
    if err != nil {
        slog.Error("Failed to get received transactions", slog.String("error", err.Error()))
        ResponseError(w, http.StatusInternalServerError, "Failed to get received transactions")
        return
    }
    var received []ReceivedTx
    for _, rt := range receivedTxs {
        received = append(received, ReceivedTx{
            FromUser: rt.FromUser,
            Amount:   rt.Amount,
        })
    }

	sentTxs, err := h.Dal.GetTransactionsSent(r.Context(), username)
    if err != nil {
        slog.Error("Failed to get sent transactions", slog.String("error", err.Error()))
        ResponseError(w, http.StatusInternalServerError, "Failed to get sent transactions")
        return
    }
    var sent []SentTx
    for _, st := range sentTxs {
        sent = append(sent, SentTx{
            ToUser: st.ToUser,
            Amount: st.Amount,
        })
    }


	resp := InfoResponse{
        Coins: user.Balance,
        Inventory: inventory,
        CoinHistory: CoinHistory{
            Received: received,
            Sent:     sent,
        },
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        slog.Error("Failed to encode info response", slog.String("error", err.Error()))
    }


}