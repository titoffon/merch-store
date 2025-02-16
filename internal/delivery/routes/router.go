package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/titoffon/merch-store/internal/db"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
)

func NewRouter(dal *db.DB) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.Handlers{
		Dal: dal,
	}

	r.Post("/api/auth", h.Auth)
	r.Get("/api/buy/{item}", h.PurchaseMerch)
	r.Post("/api/sendCoin", h.SendCoins)
	r.Get("/api/info", h.UserInfo)
	
	return r
}
