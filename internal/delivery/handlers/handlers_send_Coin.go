package handlers

import (
	"net/http"

	"github.com/titoffon/merch-store/internal/db"
)

func SendCoins(conDB *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}