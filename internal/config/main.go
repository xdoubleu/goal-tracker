//nolint:mnd //no magic number
package config

import (
	"fmt"

	"github.com/XDoubleU/essentia/pkg/config"
)

type Config struct {
	Env              string
	Port             int
	Throttle         bool
	WebURL           string
	SentryDsn        string
	SampleRate       float64
	AccessExpiry     string
	RefreshExpiry    string
	DBDsn            string
	Release          string
	GotrueProjRef    string
	GotrueAPIKey     string
	TodoistAPIKey    string
	TodoistProjectID string
	SteamAPIKey      string
	GoodreadsURL     string
}

func New() Config {
	var cfg Config

	cfg.Env = config.EnvStr("ENV", config.ProdEnv)
	cfg.Port = config.EnvInt("PORT", 8000)
	cfg.Throttle = config.EnvBool("THROTTLE", true)
	cfg.WebURL = config.EnvStr("WEB_URL", "http://localhost:3000")
	cfg.SentryDsn = config.EnvStr("SENTRY_DSN", "")
	cfg.SampleRate = config.EnvFloat("SAMPLE_RATE", 1.0)
	cfg.AccessExpiry = config.EnvStr("ACCESS_EXPIRY", "1h")
	cfg.RefreshExpiry = config.EnvStr("REFRESH_EXPIRY", "7d")
	cfg.DBDsn = config.EnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	cfg.Release = config.EnvStr("RELEASE", config.DevEnv)
	cfg.GotrueProjRef = config.EnvStr("GOTRUE_PROJ_REF", "")
	cfg.GotrueAPIKey = config.EnvStr("GOTRUE_API_KEY", "")
	cfg.TodoistAPIKey = config.EnvStr("TODOIST_API_KEY", "")
	cfg.TodoistProjectID = config.EnvStr("TODOIST_PROJECT_ID", "")
	cfg.SteamAPIKey = config.EnvStr("STEAM_API_KEY", "")
	cfg.GoodreadsURL = config.EnvStr("GOODREADS_URL", "")

	return cfg
}

func (cfg Config) String() string {
	return fmt.Sprintf(`config:
	cfg.Env: %s
	cfg.Port: %d
	cfg.Throttle: %t
	cfg.WebURL: %s
	cfg.SentryDsn: %s
	cfg.SampleRate: %f
	cfg.AccessExpiry: %s
	cfg.RefreshExpiry: %s
	cfg.DBDsn: %s
	cfg.Release: %s
	cfg.GotrueProjRef: %s
	cfg.GotrueApiKey: %s`,
		cfg.Env,
		cfg.Port,
		cfg.Throttle,
		cfg.WebURL,
		cfg.SentryDsn,
		cfg.SampleRate,
		cfg.AccessExpiry,
		cfg.RefreshExpiry,
		cfg.DBDsn,
		cfg.Release,
		cfg.GotrueProjRef,
		cfg.GotrueAPIKey,
	)
}
