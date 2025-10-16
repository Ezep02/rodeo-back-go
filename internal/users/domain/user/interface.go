package user

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, user *User) error
	UpdateUsername(ctx context.Context, new_username string, id uint) error
	UpdateAvatar(ctx context.Context, avatar string, id uint) error
}
