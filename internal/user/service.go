package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/alibekkenny/simpengine/internal/shared/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(ctx context.Context, login, email, password string) (int, error) {
	if !isValidLogin(login) {
		return 0, errors.New("invalid login format")
	}

	exists, err := s.repo.ExistsByEmailOrLogin(ctx, email, login)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, model.ErrEmailOrLoginExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &User{
		Login:    login,
		Email:    email,
		Password: string(passwordHash),
		Role:     "admin",
	}

	return s.repo.Create(ctx, user)
}

func isValidLogin(login string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
	return valid.MatchString(login)
}
