package services

import (
	"errors"
	"net/http"
	"time"

	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/xhit/go-str2duration/v2"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

type AuthService struct {
	supabaseUserID string
	client         gotrue.Client
}

func (service *AuthService) GetAllUsers() ([]models.User, error) {
	//nolint:exhaustruct //skip
	return []models.User{
		{
			ID: service.supabaseUserID,
		},
	}, nil
}

func (service *AuthService) SignInWithEmail(
	signInDto *dtos.SignInDto,
) (*string, *string, error) {
	//nolint:exhaustruct //don't need other fields
	response, err := service.client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     signInDto.Email,
		Password:  signInDto.Password,
	})
	if err != nil {
		return nil, nil, errortools.NewUnauthorizedError(
			errors.New("invalid credentials"),
		)
	}

	return &response.AccessToken, &response.RefreshToken, nil
}

func (service *AuthService) GetUser(accessToken string) (*models.User, error) {
	response, err := service.client.WithToken(accessToken).GetUser()
	if err != nil {
		return nil, err
	}

	user := models.UserFromTypesUser(response.User)

	return &user, nil
}

func (service *AuthService) SignInWithRefreshToken(
	refreshToken string,
) (*string, *string, error) {
	//nolint:exhaustruct //don't need other fields
	response, err := service.client.Token(types.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, nil, err
	}

	return &response.AccessToken, &response.RefreshToken, nil
}

func (service *AuthService) SignOut(
	accessToken string,
) (*http.Cookie, *http.Cookie, error) {
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

func (service *AuthService) GetCookieName(scope models.Scope) string {
	switch scope {
	case models.AccessScope:
		return "accessToken"
	case models.RefreshScope:
		return "refreshToken"
	default:
		panic("invalid scope")
	}
}

func (service *AuthService) CreateCookie(
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
