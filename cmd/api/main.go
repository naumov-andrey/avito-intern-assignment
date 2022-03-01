package main

import (
	"context"
	"fmt"
	"github.com/naumov-andrey/avito-intern-assignment/internal/config"
	"github.com/naumov-andrey/avito-intern-assignment/internal/handler"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository/postgres"
	"github.com/naumov-andrey/avito-intern-assignment/internal/service"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Init("configs", "main")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Config is initialized")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode)
	db, err := gorm.Open(pg.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to connect DB: %s", err)
	}
	log.Print("Database connection is established")

	r := repository.NewRepository(
		postgres.NewAccountRepositoryImpl(db),
		postgres.NewTransactionRepositoryImpl(db),
		db,
	)
	s := service.NewService(
		service.NewBalanceService(r, cfg.ExchangeAPI.AccessKey),
		service.NewTransactionService(r),
	)
	h := handler.NewHandler(s)

	server := http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: h.Init(),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Print("Server is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server is shutting down")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
