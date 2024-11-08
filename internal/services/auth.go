package services

import (
	"context"
	"errors"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"net/http"
	"time"

	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/xhit/go-str2duration/v2"
)

type AuthService struct {
	client gotrue.Client
}

func (service AuthService) SignInWithEmail(signInDto *dtos.SignInDto) (*models.User, *string, *string, error) {
	if v := signInDto.Validate(); !v.Valid() {
		return nil, nil, nil, errortools.ErrFailedValidation
	}

	response, err := service.client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     signInDto.Email,
		Password:  signInDto.Password,
	})
	if err != nil {
		return nil, nil, nil, errortools.NewUnauthorizedError(errors.New("invalid credentials"))
	}

	user := models.UserFromTypesUser(response.User)

	return &user, &response.AccessToken, &response.RefreshToken, nil
}

func (service AuthService) GetUser(accessToken string) (*models.User, error) {
	response, err := service.client.WithToken(accessToken).GetUser()
	if err != nil {
		return nil, err
	}

	user := models.UserFromTypesUser(response.User)

	return &user, nil
}

func (service AuthService) SignInWithRefreshToken(refreshToken string) (*string, *string, error) {
	response, err := service.client.Token(types.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, nil, err
	}

	return &response.AccessToken, &response.RefreshToken, nil
}

func (service AuthService) SignOut(accessToken string) (*http.Cookie, *http.Cookie, error) {
	err := service.client.WithToken(accessToken).Logout()
	if err != nil {
		return nil, nil, err
	}

	deleteAccessTokenCookie := &http.Cookie{
		Name:     service.GetCookieName(models.AccessScope),
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	}

	deleteRefreshTokenCookie := &http.Cookie{
		Name:     service.GetCookieName(models.RefreshScope),
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	}

	return deleteAccessTokenCookie, deleteRefreshTokenCookie, nil
}

func (service AuthService) GetCookieName(scope models.Scope) string {
	switch {
	case scope == models.AccessScope:
		return "accessToken"
	case scope == models.RefreshScope:
		return "refreshToken"
	default:
		panic("invalid scope")
	}
}

func (service AuthService) CreateCookie(
	ctx context.Context,
	scope models.Scope,
	token string,
	expiry string,
	secure bool,
) (*http.Cookie, error) {
	ttl, err := str2duration.ParseDuration(expiry)
	if err != nil {
		return nil, err
	}

	name := service.GetCookieName(scope)

	cookie := http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  time.Now().Add(ttl),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
	}

	return &cookie, nil
}
