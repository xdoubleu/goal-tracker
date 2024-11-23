package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/config"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/tplhelper"
)

func (app *Application) authTemplateAccess(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getCurrentUser(r)

		if user == nil {
			user = app.refreshTokens(w, r)
		}

		if user == nil {
			tplhelper.RenderWithPanic(app.tpl, w, "sign-in.html", nil)
			return
		}

		r = r.WithContext(app.contextSetUser(r.Context(), *user))
		next(w, r)
	})
}

func (app *Application) getCurrentUser(r *http.Request) *models.User {
	accessToken, err := r.Cookie("accessToken")
	if err != nil {
		return nil
	}

	user, err := app.services.Auth.GetUser(accessToken.Value)
	if err != nil {
		return nil
	}

	return user
}

func (app *Application) refreshTokens(
	w http.ResponseWriter,
	r *http.Request,
) *models.User {
	tokenCookie, err := r.Cookie("refreshToken")

	if err != nil {
		return nil
	}

	accessToken, refreshToken, err := app.services.Auth.SignInWithRefreshToken(
		tokenCookie.Value,
	)
	if err != nil {
		return nil
	}

	secure := app.config.Env == config.ProdEnv
	accessTokenCookie, err := app.services.Auth.CreateCookie(
		models.AccessScope,
		*accessToken,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		return nil
	}

	http.SetCookie(w, accessTokenCookie)

	var refreshTokenCookie *http.Cookie
	refreshTokenCookie, err = app.services.Auth.CreateCookie(
		models.RefreshScope,
		*refreshToken,
		app.config.RefreshExpiry,
		secure,
	)
	if err != nil {
		return nil
	}

	http.SetCookie(w, refreshTokenCookie)

	user, _ := app.services.Auth.GetUser(accessTokenCookie.Value)
	return user
}
