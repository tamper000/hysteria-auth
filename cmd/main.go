package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/tamper000/hysteria-auth/internal/config"
	"github.com/tamper000/hysteria-auth/internal/handler"
	"github.com/tamper000/hysteria-auth/internal/repository/sqlite"
)

const version = "0.0.1"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()

	// Create server
	webCfg := huma.DefaultConfig("Hysteria Auth", version)
	webCfg.FieldsOptionalByDefault = false
	webCfg.CreateHooks = nil

	api := humachi.New(router, webCfg)
	handlers := handler.New(db)
	registerHandlers(api, handlers)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,

		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	log.Println(fmt.Sprintf("Start web server on %s port", cfg.Port))
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func registerHandlers(api huma.API, handlers *handler.Handler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/register",
		Summary:       "Register new user",
		Description:   "Register new user",
		Tags:          []string{"API"},
		Errors:        []int{http.StatusConflict},
		DefaultStatus: http.StatusCreated,
	}, handlers.Register)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/delete",
		Summary:       "Delete user",
		Description:   "Delete user",
		Tags:          []string{"API"},
		Errors:        []int{http.StatusNotFound},
		DefaultStatus: http.StatusOK,
	}, handlers.Delete)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/auth",
		Summary:       "Hysteria auth endpoint",
		Description:   "Hysteria auth endpoint",
		Tags:          []string{"Hysteria"},
		Errors:        []int{http.StatusNotFound},
		DefaultStatus: http.StatusOK,
	}, handlers.Auth)
}
