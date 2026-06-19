package di

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"subAggregator/config"
	"subAggregator/internal/controller"
	"subAggregator/internal/domain"
	"subAggregator/internal/repository"
	"subAggregator/internal/usecase"
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
				fx.As(new(domain.SubcriptionAggregatorService)),
			),

			// HTTP-слой.
			controller.NewSubsHandler,
			NewHTTPServer,
		),
		fx.Invoke(InitDatabase),
		fx.Invoke(func(*http.Server) {}),
	)
	return module
}

func NewHTTPServer(lc fx.Lifecycle, handler *controller.SubsHandler) *http.Server {
	router := controller.NewRouter(handler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Println("starting HTTP server on", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("stopping HTTP server on", srv.Addr)
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewPostgresPool(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)

	log.Printf("Connecting to database: %s", dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("Database ping failed: %v", err)
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

func InitDatabase(pool *pgxpool.Pool, cfg *config.Config) error {
	log.Println("checking database connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("database connection established")
	if err := migrator.RunMigrations(cfg); err != nil {
		log.Printf("failed to run migrations: %s", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("database migrations completed successfully")
	return nil
}
