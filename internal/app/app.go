package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/akhmettolegen/cr-booking/internal/config"
	v1 "github.com/akhmettolegen/cr-booking/internal/handler/http/v1"
	"github.com/akhmettolegen/cr-booking/internal/repo"
	"github.com/akhmettolegen/cr-booking/internal/usecase"
	"github.com/akhmettolegen/cr-booking/pkg/logger"
	"github.com/akhmettolegen/cr-booking/pkg/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(cfg.Log.Level)

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(err)
	}
	defer pg.Close()

	reservUsecase := usecase.New(repo.New(pg))

	router := setupRouter(l, reservUsecase)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				l.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			l.Error(fmt.Errorf("failed to stop server %v", err))

			return
		}
		serverStopCtx()
	}()

	l.Info(fmt.Sprintf("listening and serving on port %s", cfg.HTTPServer.Port))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error(fmt.Errorf("failed to start server %v", err))
	}

	<-serverCtx.Done()
	l.Info("server stopped")
}

func setupRouter(l logger.Interface, r usecase.Reservation) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Mount("/swagger", httpSwagger.WrapHandler)
	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	router.Mount("/v1/reservation", v1.NewReservationHandler(l, r).Routes())

	return router
}
