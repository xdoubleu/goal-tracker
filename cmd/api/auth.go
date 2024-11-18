package main

import (
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/config"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/temptools"
)

func (app *Application) authRoutes(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s/auth/signin", prefix), app.signInHandler)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/auth/signout", prefix),
		app.authAccess(app.signOutHandler),
	)
}

func (app *Application) signInHandler(w http.ResponseWriter, r *http.Request) {
	var signInDto dtos.SignInDto

	err := temptools.ReadForm(r, &signInDto)
	if err != nil {
		temptools.RedirectWithError(w, r, "/", err)
		return
	}

	accessToken, refreshToken, err := app.services.Auth.SignInWithEmail(&signInDto)
	if err != nil {
		temptools.RedirectWithError(w, r, "/", err)
		return
	}

	secure := app.config.Env == config.ProdEnv
	accessTokenCookie, err := app.services.Auth.CreateCookie(
		models.AccessScope,
		*accessToken,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		temptools.RedirectWithError(w, r, "/", err)
		return
	}

	http.SetCookie(w, accessTokenCookie)

	if signInDto.RememberMe {
		var refreshTokenCookie *http.Cookie
		refreshTokenCookie, err = app.services.Auth.CreateCookie(
			models.RefreshScope,
			*refreshToken,
			app.config.RefreshExpiry,
			secure,
		)
		if err != nil {
			temptools.RedirectWithError(w, r, "/", err)
			return
		}

		http.SetCookie(w, refreshTokenCookie)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) signOutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, _ := r.Cookie("accessToken")
	refreshToken, _ := r.Cookie("refreshToken")

	deleteAccessTokenCookie, deleteRefreshTokenCookie, err := app.services.Auth.SignOut(
		accessToken.Value,
	)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, deleteAccessTokenCookie)

	if refreshToken != nil {
		http.SetCookie(w, deleteRefreshTokenCookie)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
