package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "nds-go-starter/docs"
	kafkaAdapter "nds-go-starter/internal/adapters/kafka"
	"nds-go-starter/internal/core/auth"
	"nds-go-starter/internal/core/repository"
	"nds-go-starter/internal/core/worker"
	authFeature "nds-go-starter/internal/features/auth"
	"nds-go-starter/internal/features/notifications"
	"nds-go-starter/internal/features/orders"
	"nds-go-starter/internal/features/products"
	"nds-go-starter/internal/features/workout"
	appMiddleware "nds-go-starter/internal/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(appMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Configurable CORS
	if app.config.corsEnabled {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   app.config.corsOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	var db repository.DBTX = app.db
	if app.config.dbLog {
		db = repository.NewLogDB(app.db)
	}

	authFeature.Register(r, db, app.jwtManager)

	// Kafka Initialization
	kafkaAddr := "localhost:9092"
	orderCreatedWriter := kafkaAdapter.NewWriter(kafkaAddr, "order-created")
	orderCreatedReader := kafkaAdapter.NewReader(kafkaAddr, "order-created", "notification-group")
	productCreatedWriter := kafkaAdapter.NewWriter(kafkaAddr, "product-created")
	// productCreatedReader := kafkaAdapter.NewReader(kafkaAddr, "product-created", "notification-group")

	// Register Workers
	authRepo := authFeature.NewRepository(repository.New(db))
	cleanupWorker := authFeature.NewSessionCleanupWorker(authRepo, 1*time.Hour)
	notificationWorker := notifications.NewNotificationWorker(orderCreatedReader)
	// productNotificationWorker := notifications.NewNotificationWorker(productCreatedReader)

	app.addWorker(cleanupWorker)
	app.addWorker(notificationWorker)

	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.Auth(app.jwtManager, app.config.authEnabled))
		products.Register(r, db, productCreatedWriter)
		orders.Register(r, db, orderCreatedWriter)

		workout.Register(r, db)
	})

	return r
}

func (app *application) run(ctx context.Context, h http.Handler) error {
	for _, w := range app.workers {
		go w.Start(ctx)
	}

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

type application struct {
	config     config
	db         *sql.DB
	jwtManager *auth.JWTManager
	workers    []worker.Worker
}

func (app *application) addWorker(w worker.Worker) {
	app.workers = append(app.workers, w)
}

type config struct {
	addr        string
	db          dbConfig
	dbLog       bool
	authEnabled bool
	jwtSecret   string
	logFile     string
	corsEnabled bool
	corsOrigins []string
}

type dbConfig struct {
	dsn string
}
