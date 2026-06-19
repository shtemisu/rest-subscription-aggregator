package di

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"subAggregator/config"
	"subAggregator/internal/controller"
	"subAggregator/internal/domain"
	"subAggregator/internal/repository"
	"subAggregator/internal/usecase"
	"subAggregator/pkg/logger"
	"subAggregator/pkg/migrator"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func Module() fx.Option {
	module := fx.Module("subs-aggregator",
		fx.Provide(
			config.NewConfig,
			NewPostgresPool,

			fx.Annotate(
				repository.NewSubsRepo,
				fx.As(new(domain.SubsInfoRepository)),
			),

			fx.Annotate(
				usecase.NewSubAggregatorService,
				fx.As(new(domain.SubscriptionAggregatorService)),
			),
			logger.New,
			controller.NewSubsHandler,
			NewHTTPServer,
		),
		fx.Invoke(InitDatabase),
		fx.Invoke(func(*http.Server) {}),
	)
	return module
}

func NewHTTPServer(lc fx.Lifecycle, handler *controller.SubsHandler, l *slog.Logger, shutdowner fx.Shutdowner) *http.Server {
	router := controller.NewRouter(handler)
	h := controller.LoggingMiddleware(l)(router)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			l.Info("starting HTTP server", "address", srv.Addr)
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					l.Error("http server stopped unexpectedly", "error", err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("stopping HTTP server", "address", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewPostgresPool(lc fx.Lifecycle, cfg *config.Config, l *slog.Logger) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		l.Error("Database ping failed", "error", err)
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool, nil
}

func InitDatabase(pool *pgxpool.Pool, cfg *config.Config, l *slog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	l.Info("database connection established")
	if err := migrator.RunMigrations(cfg); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	l.Info("database migrations completed successfully")
	return nil
}
