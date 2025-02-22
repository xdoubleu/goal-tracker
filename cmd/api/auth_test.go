package main

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/dtos"
)

func TestSignInHandler(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPost,
		"/api/auth/signin",
	)

	signInDto := dtos.SignInDto{
		Email:      "valid@example.com",
		Password:   "password",
		RememberMe: true,
	}

	tReq.SetFollowRedirect(false)

	tReq.SetContentType(test.FormContentType)
	tReq.SetData(signInDto)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusSeeOther, rs.StatusCode)
}

func TestSignOutHandler(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/api/auth/signout",
	)

	tReq.SetFollowRedirect(false)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusSeeOther, rs.StatusCode)
}
