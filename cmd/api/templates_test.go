package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

func TestSignIn(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/",
	)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestRefreshTokens(t *testing.T) {
	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/",
	)

	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestRoot(t *testing.T) {
	err := testApp.services.Goals.ImportStatesFromTodoist(context.Background(), userID)
	if err != nil {
		panic(err)
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
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
		http.MethodGet,
		"/edit/123",
	)
	tReq.AddCookie(&accessToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestGoalProgressGraph(t *testing.T) {
	err := testApp.services.Goals.ImportGoalsFromTodoist(
		context.Background(),
		testApp.config.SupabaseUserID,
	)
	if err != nil {
		panic(err)
	}

	val := int64(50)
	err = testApp.services.Goals.LinkGoal(
		context.Background(),
		goalID,
		userID,
		&dtos.LinkGoalDto{
			TypeID:      models.SteamCompletionRate.ID,
			TargetValue: &val,
			Tag:         nil,
		},
	)
	if err != nil {
		panic(err)
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/goals/123",
	)
	tReq.AddCookie(&accessToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestGoalProgressList(t *testing.T) {
	err := testApp.services.Goals.ImportGoalsFromTodoist(
		context.Background(),
		testApp.config.SupabaseUserID,
	)
	if err != nil {
		panic(err)
	}

	val := int64(50)
	valStr := "fiction"
	err = testApp.services.Goals.LinkGoal(
		context.Background(),
		goalID,
		userID,
		&dtos.LinkGoalDto{
			TypeID:      models.BooksFromSpecificTag.ID,
			TargetValue: &val,
			Tag:         &valStr,
		},
	)
	if err != nil {
		panic(err)
	}

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		"/goals/123",
	)
	tReq.AddCookie(&accessToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
