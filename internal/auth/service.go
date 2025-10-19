package auth

import (
	"context"
	"errors"
	"fmt"

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
	foundUser, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return "", fmt.Errorf("%w: user with such login and password not found", model.ErrInvalidCredentials)
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%w: user with such login and password not found", model.ErrInvalidCredentials)
	}

	return GenerateJWT(foundUser.ID, foundUser.Login, foundUser.Role)
}

func (s *AuthService) Exists(ctx context.Context, id int64) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}
