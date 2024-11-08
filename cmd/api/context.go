package main

import (
	"context"
	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"

	"github.com/getsentry/sentry-go"
)

func (app *Application) contextSetUser(
	ctx context.Context,
	user models.User,
) context.Context {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		//nolint:exhaustruct //other fields are optional
		hub.Scope().SetUser(sentry.User{
			ID:    user.ID,
			Email: user.Email,
		})
	}

	return context.WithValue(ctx, constants.UserContextKey, user)
}
