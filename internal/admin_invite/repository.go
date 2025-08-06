package admininvite

import "context"

type AdminInviteRepository interface {
	Create(ctx context.Context, invite AdminInvite) (int, error)
	GetInviteByToken(ctx context.Context, token string) (*AdminInvite, error)
	MarkInviteAsUsedByToken(ctx context.Context, token string, usedById int) error
}
