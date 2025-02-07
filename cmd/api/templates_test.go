package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		"/",
	)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestRefreshTokens(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		"/",
	)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestRoot(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		"/",
	)
	tReq.AddCookie(&accessToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestLink(t *testing.T) {
	err := testApp.services.Goals.ImportGoalsFromTodoist(
		context.Background(),
		testApp.config.SupabaseUserID,
	)
	if err != nil {
		panic(err)
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		"/link/123",
	)
	tReq.AddCookie(&accessToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
