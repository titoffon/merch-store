package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/titoffon/merch-store/internal/db"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
)

func NewRouter(connectionDB *db.DB) *chi.Mux {
	r := chi.NewRouter()


	r.Post("/api/auth", handlers.Auth(connectionDB))

	r.Route("/api", func(r chi.Router) {
		r.Get("/info", handlers.GetUserInfo(connectionDB))
		r.Get("/me/purchases", handlers.GetUserPurchases(connectionDB))
		r.Get("/me/transactions", handlers.GetUserTransactions(connectionDB))
	})

	r.Route("/merch", func(r chi.Router) {
		r.Post("/purchase", handlers.PurchaseMerch(connectionDB))
		r.Post("/coins/transfer", handlers.TransferCoins(connectionDB))
	})

	return r
}
