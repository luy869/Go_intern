package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"yatter-backend-go/app/config"
	"yatter-backend-go/app/domain/service"
	"yatter-backend-go/app/infra"
	"yatter-backend-go/app/infra/query"
	"yatter-backend-go/app/infra/transaction"
	api_auth "yatter-backend-go/app/ui/api/auth"
	"yatter-backend-go/app/ui/api/health"
	api_timeline "yatter-backend-go/app/ui/api/timeline"
	api_user "yatter-backend-go/app/ui/api/user"
	api_yweet "yatter-backend-go/app/ui/api/yweet"
	"yatter-backend-go/app/usecase/auth"
	"yatter-backend-go/app/usecase/user"
	"yatter-backend-go/app/usecase/yweet"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

const (
	requestTimeout    = 60 * time.Second
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func Run(db *sqlx.DB) error {
	addr := ":" + strconv.Itoa(config.Port())

	transactor := transaction.NewTransactor(db)

	userRepo := infra.NewUserRepository()
	yweetRepo := infra.NewYweetRepository()

	usernameUniqueChecker := service.NewUsernameUniqueChecker(userRepo)

	authQueryService := query.NewAuthQueryService(db)
	timelineQueryService := query.NewTimelineQueryService(db)
	yweetQueryService := query.NewYweetQueryService(db)
	userQueryService := query.NewUserQueryService(db)

	userCreateUseCase := user.NewUserCreateUseCase(
		userRepo,
		usernameUniqueChecker,
		transactor,
	)
	loginUseCase := auth.NewLoginUseCase(authQueryService)
	createYweetUseCase := yweet.NewCreateYweetUseCase(userRepo, yweetRepo, transactor)
	updateCredentialUseCase := user.NewUpdateCredentialUseCase(userRepo, transactor)

	userHandler := api_user.NewUserHandler(userCreateUseCase, userQueryService, updateCredentialUseCase)
	authHandler := api_auth.NewAuthHandler(loginUseCase)
	yweetHandler := api_yweet.NewYweetHandler(createYweetUseCase, yweetQueryService)
	timelineHandler := api_timeline.NewTimelineHandler(timelineQueryService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(newCORS().Handler)

	r.Use(middleware.Timeout(requestTimeout))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.SignUp)
			r.Get("/{username}", userHandler.FindByUsername)
			r.Post("/update_credentials", userHandler.UpdateUserCredential)
		})

		r.Route("/health", func(r chi.Router) {
			r.Get("/", health.Check)
		})

		r.Route("/yweets", func(r chi.Router) {
			r.Post("/", yweetHandler.Create)
			r.Get("/{id}", yweetHandler.Find)
		})

		r.Route("/timelines", func(r chi.Router) {
			r.Get("/public", timelineHandler.Public)
		})
	})

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	l, err := net.Listen("tcp", addr)
	slog.Info("Serve on 127.0.0.1", "addr", addr)
	if err != nil {
		slog.Error("failed to listen", "err", err)
	}

	go func() {
		if err = srv.Serve(l); err != nil {
			slog.Error("failed to serve", "err", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "err", err)
	}

	return nil
}

func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	})
}
