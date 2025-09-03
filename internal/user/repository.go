package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	FindById(ctx context.Context, id int64) (*User, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	ExistsByEmailOrLogin(ctx context.Context, email, login string) (bool, error)
	ExistsByID(ctx context.Context, id int64) (bool, error)
}
