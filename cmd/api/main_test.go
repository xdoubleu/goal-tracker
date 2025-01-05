package main

import (
	"context"
	"os"
	"testing"
	"time"

	configtools "github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"

	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/mocks"
)

type TestEnv struct {
	ctx context.Context
	tx  *postgres.PgxSyncTx
	app *Application
}

var mainTx *postgres.PgxSyncTx //nolint:gochecknoglobals //needed for tests
var cfg config.Config          //nolint:gochecknoglobals //needed for tests
var mainTestApp *Application   //nolint:gochecknoglobals //needed for tests
var testCtx context.Context    //nolint:gochecknoglobals //needed for tests

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false
	cfg.SupabaseUserID = "4001e9cf-3fbe-4b09-863f-bd1654cfbf76"

	postgresDB, err := postgres.Connect(
		logging.NewNopLogger(),
		cfg.DBDsn,
		25,
		"15m",
		5,
		15*time.Second,
		30*time.Second,
	)
	if err != nil {
		panic(err)
	}

	ApplyMigrations(logging.NewNopLogger(), postgresDB)

	mainTx = postgres.CreatePgxSyncTx(context.Background(), postgresDB)

	clients := Clients{
		Supabase:  mocks.NewMockedGoTrueClient(),
		Todoist:   mocks.NewMockTodoistClient(),
		Steam:     mocks.NewMockSteamClient(),
		Goodreads: mocks.NewMockGoodreadsClient(),
	}

	mainTestApp = NewApp(
		logging.NewNopLogger(),
		cfg,
		mainTx,
		clients,
	)
	testCtx = context.Background()

	code := m.Run()

	err = mainTx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func setup(_ *testing.T) (*TestEnv, *Application) {
	tx := postgres.CreatePgxSyncTx(context.Background(), mainTx)

	testApp := *mainTestApp
	testApp.setDB(tx)

	testEnv := &TestEnv{
		ctx: testCtx,
		tx:  tx,
		app: &testApp,
	}

	return testEnv, &testApp
}

func (env *TestEnv) teardown() {
	err := env.tx.Rollback(context.Background())
	if err != nil {
		panic(err)
	}

	env.app.ctxCancel()
}
