package main

import (
	"time"

	"github.com/mayankpatidar275/go-social/internal/db"
	"github.com/mayankpatidar275/go-social/internal/env"
	"github.com/mayankpatidar275/go-social/internal/mailer"
	"github.com/mayankpatidar275/go-social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	// This file will mostly have the configurations
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/go-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// Note: passing the database connection to storage layer which implements all the details
	// Our handlers will receive the storage
	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	// Note:
	// Type: application is the blueprint (no memory used until instantiated).
	// Value: application{} creates the actual object in memory.
	// Pointer: &application{} points to the memory location of the created object.

	app := &applicaion{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
