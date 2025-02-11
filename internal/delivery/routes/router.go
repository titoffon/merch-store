package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
)

func NewRouter(pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

r.Post("/auth/register", handlers.Register(pool))
	r.Post("/auth/login", handlers.Login(pool))

	r.Route("/users", func(r chi.Router) {
		r.Get("/me", handlers.GetUserInfo(pool))
		r.Get("/me/purchases", handlers.GetUserPurchases(pool))
		r.Get("/me/transactions", handlers.GetUserTransactions(pool))
	})

	r.Route("/merch", func(r chi.Router) {
		r.Post("/purchase", handlers.PurchaseMerch(pool))
		r.Post("/coins/transfer", handlers.TransferCoins(pool))
	})

	return r
}
