package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

/* returns http.Handler*/
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/api/v1/authentication/signup", app.Signup)
	mux.Post("/api/v1/authentication/login", app.Login)
	mux.Get("/api/v1/authentication/all-users", app.GetAllUsers)

	//POS Mainnet & Testnet: Missed Checkpoint
	mux.Get("/api/v1/pos/mainnet/mainnet-missed-checkpoint", app.MainnetMissedCheckpoint)
	mux.Get("/api/v1/pos/testnet/mainnet-missed-checkpoint", app.TestnetMissedCheckpoint)

	//POS Mainnet &Testnet: Heimdall Block Height
	mux.Get("/api/v1/pos/mainnet/heimdal-block-height", app.MainnetHeimdalBlockHeight)
	mux.Get("/api/v1/pos/testnet/heimdal-block-height", app.TestnetHeimdalBlockHeight)

	//POS Mainnet &Testnet: Bor Latest Block Detail
	mux.Get("/api/v1/pos/mainnet/bor-latest-block-details", app.MainnetBorLatestBlockDetails)
	mux.Get("/api/v1/pos/testnet/bor-latest-block-details", app.TestnetBorLatestBlockDetails)

	// mux.Get("/api/v1/authentication/get-me", app.GetMe)
	// mux.Get("/api/v1/authentication/verify-token", app.VerifyToken)
	// mux.Post("/api/v1/authentication/log-out", app.Logout)

	return mux
}
