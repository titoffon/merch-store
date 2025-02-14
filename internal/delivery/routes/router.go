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
	//r.Get("/api/sendCoin", h.SendCoins)

	/*r.Route("/api", func(r chi.Router) {
		r.Get("/info", handlers.GetUserInfo(dal))
		r.Get("/me/purchases", handlers.GetUserPurchases(dal))
		r.Get("/me/transactions", handlers.GetUserTransactions(dal))
	})

	r.Route("/merch", func(r chi.Router) {
		r.Post("/purchase", handlers.PurchaseMerch(dal))
		r.Post("/coins/transfer", handlers.TransferCoins(dal))
	})*/

	return r
}
