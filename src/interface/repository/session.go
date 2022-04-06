package repository

import (
	"auth-project/src/domain/model"
	"context"
	"github.com/uptrace/bun"
)

type sessionRepository struct {
	db *bun.DB
}

type SessionRepository interface {
	InsertSession(ctx context.Context, ses *model.Session) error
	UpdateSession(ctx context.Context, ses *model.Session) error
}

func NewSessionRepository(db *bun.DB) SessionRepository {
	return &sessionRepository{db}
}

func (sr *sessionRepository) InsertSession(ctx context.Context, ses *model.Session) error {
	_, err := sr.db.NewInsert().Model(ses).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepository) UpdateSession(ctx context.Context, ses *model.Session) error {
	_, err := sr.db.NewUpdate().Model(ses).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
