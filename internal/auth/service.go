package auth

import (
	"context"

	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/alibekkenny/simpengine/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo user.UserRepository
}

func NewAuthService(r user.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", model.ErrNoRecord
	}

	jwtToken, err := GenerateJWT(user.ID, user.Login)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
