package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
)

func TestEditGoalHandler(t *testing.T) {
	//nolint:exhaustruct,errcheck //other fields are optional
	defer testApp.repositories.Goals.Delete(
		context.Background(),
		&models.Goal{ID: goalID},
		userID,
	)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodPost,
		fmt.Sprintf("/api/goals/%s/edit", goalID),
	)

	tReq.SetFollowRedirect(false)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	targetValue := int64(50)
	tag := ""

	tReq.SetContentType(test.FormContentType)
	tReq.SetData(dtos.LinkGoalDto{
		TypeID:      models.SteamCompletionRate.ID,
		TargetValue: &targetValue,
		Tag:         &tag,
	})

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusSeeOther, rs.StatusCode)
}

func TestUnlinkGoalHandler(t *testing.T) {
	_, err := testApp.repositories.Goals.Upsert(
		context.Background(),
		goalID,
		userID,
		"Goal",
		"1",
		nil,
		1,
	)
	if err != nil {
		panic(err)
	}
	//nolint:exhaustruct,errcheck //other fields are optional
	defer testApp.repositories.Goals.Delete(
		context.Background(),
		&models.Goal{ID: goalID},
		userID,
	)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		fmt.Sprintf("/api/goals/%s/unlink", goalID),
	)

	tReq.SetFollowRedirect(false)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusSeeOther, rs.StatusCode)
}

func TestCompleteGoalHandler(t *testing.T) {
	_, err := testApp.repositories.Goals.Upsert(
		context.Background(),
		goalID,
		userID,
		"Goal",
		"1",
		nil,
		1,
	)
	if err != nil {
		panic(err)
	}
	//nolint:exhaustruct,errcheck //other fields are optional
	defer testApp.repositories.Goals.Delete(
		context.Background(),
		&models.Goal{ID: goalID},
		userID,
	)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		http.MethodGet,
		fmt.Sprintf("/api/goals/%s/complete", goalID),
	)

	tReq.SetFollowRedirect(false)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusSeeOther, rs.StatusCode)
}
