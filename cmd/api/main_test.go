package main

import (
	"context"
	"goal-tracker/api/internal/config"
	"goal-tracker/api/internal/mocks"
	"goal-tracker/api/pkg/todoist"
	"os"
	"testing"
	"time"

	configtools "github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/supabase-community/gotrue-go"
)

type TestEnv struct {
	ctx context.Context
	tx  *postgres.PgxSyncTx
	app *Application
}

var mainTx *postgres.PgxSyncTx   //nolint:gochecknoglobals //needed for tests
var cfg config.Config            //nolint:gochecknoglobals //needed for tests
var mainTestApp *Application     //nolint:gochecknoglobals //needed for tests
var testCtx context.Context      //nolint:gochecknoglobals //needed for tests
var supabaseClient gotrue.Client //nolint:gochecknoglobals //needed for tests
var todoistClient todoist.Client //nolint:gochecknoglobals //needed for tests

func TestMain(m *testing.M) {
	var err error

	cfg = config.New()
	cfg.Env = configtools.TestEnv
	cfg.Throttle = false

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
	supabaseClient = mocks.NewMockedGoTrueClient()
	todoistClient = mocks.NewMockTodoistClient()
	mainTestApp = NewApp(logging.NewNopLogger(), cfg, mainTx, supabaseClient, todoistClient)
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
	testApp.setDB(tx, supabaseClient, todoistClient)

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
