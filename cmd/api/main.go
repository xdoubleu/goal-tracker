package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	sentrytools "github.com/XDoubleU/essentia/pkg/sentry"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/supabase-community/gotrue-go"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/jobs"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/internal/services"
	"goal-tracker/api/internal/temptools"
	"goal-tracker/api/pkg/goodreads"
	"goal-tracker/api/pkg/steam"
	"goal-tracker/api/pkg/todoist"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

//go:embed templates/html/**/*html
var htmlTemplates embed.FS

//go:embed images/**
var images embed.FS

type Application struct {
	logger    *slog.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc
	db        postgres.DB
	config    config.Config
	images    embed.FS
	clients   Clients
	services  services.Services
	tpl       *template.Template
	jobQueue  *temptools.JobQueue
}

type Clients struct {
	Supabase  gotrue.Client
	Steam     steam.Client
	Todoist   todoist.Client
	Goodreads goodreads.Client
}

//	@title			goal-tracker API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	logger := slog.New(sentrytools.NewLogHandler(cfg.Env,
		slog.NewTextHandler(os.Stdout, nil)))
	db, err := postgres.Connect(
		logger,
		cfg.DBDsn,
		25, //nolint:mnd //no magic number
		"15m",
		60,             //nolint:mnd //no magic number
		10*time.Second, //nolint:mnd //no magic number
		5*time.Minute,  //nolint:mnd //no magic number
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ApplyMigrations(logger, db)

	clients := Clients{
		Supabase: gotrue.New(
			cfg.SupabaseProjRef,
			cfg.SupabaseAPIKey,
		),
		Todoist:   todoist.New(cfg.TodoistAPIKey),
		Steam:     steam.New(logger, cfg.SteamAPIKey),
		Goodreads: goodreads.New(logger),
	}

	app := NewApp(logger, cfg, db, clients)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:mnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:mnd //no magic number
	}
	err = httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}

func NewApp(
	logger *slog.Logger,
	cfg config.Config,
	db postgres.DB,
	clients Clients,
) *Application {
	tpl := template.Must(template.ParseFS(htmlTemplates, "templates/html/**/*.html"))

	//nolint:mnd //no magic number
	jobQueue := temptools.NewJobQueue(*logger, 100)

	//nolint:exhaustruct //other fields are optional
	app := &Application{
		logger:   logger,
		clients:  clients,
		config:   cfg,
		images:   images,
		tpl:      tpl,
		jobQueue: &jobQueue,
	}

	app.setContext()
	app.setDB(db)

	err := app.jobQueue.Push(
		jobs.NewTodoistJob(app.services.Auth, app.services.Goals),
		app.services.WebSocket.UpdateState,
	)
	if err != nil {
		panic(err)
	}

	err = app.jobQueue.Push(
		jobs.NewGoodreadsJob(
			app.services.Auth,
			app.services.Goodreads,
			app.services.Goals,
		),
		app.services.WebSocket.UpdateState,
	)
	if err != nil {
		panic(err)
	}

	err = app.jobQueue.Push(
		jobs.NewSteamJob(app.services.Auth, app.services.Steam, app.services.Goals),
		app.services.WebSocket.UpdateState,
	)
	if err != nil {
		panic(err)
	}

	app.services.WebSocket.RegisterTopics(app.jobQueue.FetchRecurringJobIDs())

	return app
}

func (app *Application) setDB(
	db postgres.DB,
) {
	// make sure previous app is cancelled internally
	app.ctxCancel()

	app.setContext()

	spandb := postgres.NewSpanDB(db)

	app.db = spandb
	app.services = services.New(
		*app.logger,
		app.config,
		app.jobQueue,
		repositories.New(app.db),
		app.clients.Supabase,
		app.clients.Todoist,
		app.clients.Steam,
		app.clients.Goodreads,
	)
}

func (app *Application) setContext() {
	ctx, cancel := context.WithCancel(context.Background())
	app.ctx = ctx
	app.ctxCancel = cancel
}

func ApplyMigrations(logger *slog.Logger, db *pgxpool.Pool) {
	migrationsDB := stdlib.OpenDBFromPool(db)

	goose.SetLogger(slog.NewLogLogger(logger.Handler(), slog.LevelInfo))

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		panic(err)
	}

	if err := goose.Up(migrationsDB, "migrations"); err != nil {
		panic(err)
	}
}
