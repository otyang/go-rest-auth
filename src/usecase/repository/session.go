package repository

import (
	"auth-project/src/domain/model"
	"context"
)

type SessionRepository interface {
	InsertSession(ctx context.Context, ses *model.Session) error
	UpdateSession(ctx context.Context, ses *model.Session) error
}
