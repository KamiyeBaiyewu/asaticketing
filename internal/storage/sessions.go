package storage

import (
	"context"

	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

func (s *storage) SaveRefreshToken(ctx context.Context, session *model.Session) error {
	return s.db.SaveRefreshToken(ctx, session)
}
func (s *storage) GetSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	return s.db.GetSession(ctx, session)
}
