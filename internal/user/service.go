package user

import (
	"context"
	"fmt"
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

func (s *UserService) Register(ctx context.Context, login, email, password string) (int64, error) {
	if !isValidLogin(login) {
		return 0, fmt.Errorf("%w: invalid login format", model.ErrInvalidBody)
	}

	exists, err := s.repo.ExistsByEmailOrLogin(ctx, email, login)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("%w: user with such login or email already exists", model.ErrUniqueViolation)
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
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
}

func isValidLogin(login string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
	return valid.MatchString(login)
}
