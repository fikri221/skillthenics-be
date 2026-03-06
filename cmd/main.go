package main

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"nds-go-starter/internal/core/auth"
	"nds-go-starter/internal/env"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
)

// @title NDS Go Starter API
// @version 1.0
// @description This is a sample server for a Go project with Clean Architecture.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8788
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your-jwt-token>" to authenticate.

func main() {
	env.Load()
	//ctx := context.Background()

	cfg := config{
		addr: ":8788",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "root:passwordnds@tcp(localhost:3306)/db_nds_go_starter?parseTime=true"),
		},
		dbLog:       env.GetBool("DB_LOG", false),
		authEnabled: env.GetBool("AUTH_ENABLED", false),
		jwtSecret:   env.GetString("JWT_SECRET", "default-secret"),
		logFile:     env.GetString("LOG_FILE_PATH", "logs/app.log"),
		corsEnabled: env.GetBool("CORS_ENABLED", false),
		corsOrigins: strings.Split(env.GetString("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}

	jwtManager := auth.NewJWTManager(cfg.jwtSecret, 15*time.Minute)

	// Enterprise Logger: Multi-Writer (Console + File) with Rotation
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.logFile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
	logger := slog.New(slog.NewJSONHandler(multiWriter, nil))
	slog.SetDefault(logger)

	logger.Info("Database connection string", "dsn", cfg.db.dsn)

	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api := application{
		config:     cfg,
		db:         db,
		jwtManager: jwtManager,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := api.run(ctx, api.mount()); err != nil {
		slog.Error("Server has failed to start", "err", err)
		os.Exit(1)
	}

}
