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

func TestLinkGoalHandler(t *testing.T) {
	//nolint:exhaustruct,errcheck //other fields are optional
	defer testApp.repositories.Goals.Delete(
		context.Background(),
		&models.Goal{ID: goalID},
		userID,
	)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodPost,
		fmt.Sprintf("/api/goals/%s/link", goalID),
	)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	targetValue := int64(50)
	tag := ""

	tReq.SetData(dtos.LinkGoalDto{
		TypeID:      models.SteamCompletionRate.ID,
		TargetValue: &targetValue,
		Tag:         &tag,
	})

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}

func TestUnlinkGoalHandler(t *testing.T) {
	//nolint:exhaustruct,errcheck //other fields are optional
	defer testApp.repositories.Goals.Delete(
		context.Background(),
		&models.Goal{ID: goalID},
		userID,
	)

	tReq := test.CreateRequestTester(
		testApp.routes(),
		test.FormContentType,
		http.MethodGet,
		fmt.Sprintf("/api/goals/%s/unlink", goalID),
	)

	tReq.AddCookie(&accessToken)
	tReq.AddCookie(&refreshToken)

	rs := tReq.Do(t)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
