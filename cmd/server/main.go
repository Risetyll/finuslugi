package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/risetyll/finuslugi/internal/handlers"
	"github.com/risetyll/finuslugi/internal/logger"
	lu "github.com/risetyll/finuslugi/internal/logger/utils"
	"github.com/risetyll/finuslugi/internal/storage/postgres"
	"github.com/risetyll/finuslugi/internal/usecase"
)

func setupRoutes(uc *usecase.DatabaseUsecase, logger *slog.Logger) *mux.Router {
	mux := mux.NewRouter()
	routes := handlers.New(logger)

	mux.HandleFunc("/", routes.GetRoutesInfoHandler()).Methods("GET")
	mux.HandleFunc("/materials", routes.CreateMaterialHandler(uc)).Methods("POST")
	mux.HandleFunc("/materials/{id:[0-9]+}", routes.GetMaterialByIdHandler(uc)).Methods("GET")
	mux.HandleFunc("/materials", routes.UpdateMaterialHandler(uc)).Methods("PUT")
	mux.HandleFunc("/materials", routes.GetMaterialsHandler(uc)).Methods("GET")

	return mux
}

func main() {
	logger := logger.New()
	logger.Info("logger initialization successful")

	logger.Info("connecting to database")
	postgresConnector := postgres.PostgresConnector{}
	postgres, err := postgres.New(postgresConnector, logger)
	if err != nil {
		logger.Error("connecting to database", lu.Error(err))
	}
	err = postgres.Init()
	if err != nil {
		logger.Error("database initialization error", lu.Error(err))
	}
	usecase := usecase.New(postgres)

	router := setupRoutes(usecase, logger)
	server := http.Server{
		Addr:    ":8443",
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS12,
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("starting server", slog.Any("port", server.Addr))
	if err := server.ListenAndServeTLS(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE")); err != nil {
		logger.Error("starting server", lu.Error(err))
	}
}
