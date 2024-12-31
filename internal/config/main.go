//nolint:mnd //no magic number
package config

import (
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
	SteamUserID      string
	GoodreadsURL     string
}

func New() Config {
	var cfg Config

	cfg.Env = config.EnvStr("ENV", config.ProdEnv)
	cfg.Port = config.EnvInt("PORT", 8000)
	cfg.Throttle = config.EnvBool("THROTTLE", true)
	cfg.WebURL = config.EnvStr("WEB_URL", "http://localhost:8000")
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
	cfg.SteamUserID = config.EnvStr("STEAM_USER_ID", "")

	cfg.GoodreadsURL = config.EnvStr("GOODREADS_URL", "")

	return cfg
}
