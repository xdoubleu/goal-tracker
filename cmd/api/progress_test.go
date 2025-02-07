package main

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRefreshProgressHandler(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		"/api/progress/0/refresh",
	)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
