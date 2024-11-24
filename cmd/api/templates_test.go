package main

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/",
	)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestRoot(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testApp.setDB(testEnv.tx, supabaseClient, todoistClient, steamClient)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/",
	)
	tReq.AddCookie(&http.Cookie{Name: "accessToken", Value: "thisisavaliduser"})

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestLink(t *testing.T) {
	testEnv, testApp := setup(t)
	defer testEnv.teardown()

	testApp.setDB(testEnv.tx, supabaseClient, todoistClient, steamClient)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/link/123",
	)
	tReq.AddCookie(&http.Cookie{Name: "accessToken", Value: "thisisavaliduser"})

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
