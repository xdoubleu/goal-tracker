package main

import (
	"errors"
	"net/http"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/config"
	errortools "github.com/XDoubleU/essentia/pkg/errors"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

func (app *Application) authRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /auth/current-user",
		app.authAccess(app.getCurrentUserHandler),
	)
	mux.HandleFunc("POST /auth/signin", app.signInHandler)
	mux.HandleFunc(
		"GET /auth/signout",
		app.authAccess(app.signOutHandler),
	)
	mux.HandleFunc(
		"GET /auth/refresh",
		app.refreshHandler,
	)
}

// @Summary	Get current user
// @Tags		auth
// @Success	200			{object}	User
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/auth/current-user [get].
func (app *Application) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, _ := r.Cookie("accessToken")

	user, err := app.services.Auth.GetUser(accessToken.Value)
	if err != nil {
		httptools.HandleError(w, r, err, nil)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Sign in a user
// @Tags		auth
// @Param		signInDto	body		SignInDto	true	"SignInDto"
// @Success	200			{object}	User
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/auth/signin [post].
func (app *Application) signInHandler(w http.ResponseWriter, r *http.Request) {
	var signInDto *dtos.SignInDto

	err := httptools.ReadJSON(r.Body, &signInDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user, accessToken, refreshToken, err := app.services.Auth.SignInWithEmail(signInDto)
	if err != nil {
		httptools.HandleError(w, r, err, signInDto.ValidationErrors)
		return
	}

	secure := app.config.Env == config.ProdEnv
	accessTokenCookie, err := app.services.Auth.CreateCookie(
		r.Context(),
		models.AccessScope,
		*accessToken,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, accessTokenCookie)

	if signInDto.RememberMe {
		var refreshTokenCookie *http.Cookie
		refreshTokenCookie, err = app.services.Auth.CreateCookie(
			r.Context(),
			models.RefreshScope,
			*refreshToken,
			app.config.RefreshExpiry,
			secure,
		)
		if err != nil {
			httptools.ServerErrorResponse(w, r, err)
			return
		}

		http.SetCookie(w, refreshTokenCookie)
	}

	err = httptools.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Sign out a user
// @Tags		auth
// @Success	200	{object}	nil
// @Failure	401	{object}	ErrorDto
// @Router		/auth/signout [get].
func (app *Application) signOutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, _ := r.Cookie("accessToken")
	refreshToken, _ := r.Cookie("refreshToken")

	deleteAccessTokenCookie, deleteRefreshTokenCookie, err := app.services.Auth.SignOut(accessToken.Value)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, deleteAccessTokenCookie)

	if refreshToken == nil {
		return
	}

	http.SetCookie(w, deleteRefreshTokenCookie)
}

// @Summary	Refresh access token
// @Tags		auth
// @Success	200	{object}	nil
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/auth/refresh [get].
func (app *Application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("refreshToken")

	if err != nil {
		httptools.UnauthorizedResponse(w, r,
			errortools.NewUnauthorizedError(errors.New("no token in cookies")))
		return
	}

	accessToken, refreshToken, err := app.services.Auth.SignInWithRefreshToken(tokenCookie.Value)
	if err != nil {
		httptools.HandleError(w, r, err, nil)
		return
	}

	secure := app.config.Env == config.ProdEnv
	accessTokenCookie, err := app.services.Auth.CreateCookie(
		r.Context(),
		models.AccessScope,
		*accessToken,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, accessTokenCookie)

	var refreshTokenCookie *http.Cookie
	refreshTokenCookie, err = app.services.Auth.CreateCookie(
		r.Context(),
		models.RefreshScope,
		*refreshToken,
		app.config.RefreshExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, refreshTokenCookie)
}
