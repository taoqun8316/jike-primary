package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"jike/internal/domain"
	"jike/internal/repository"
)

var (
	ErrDuplicateEmail      = repository.ErrDuplicateEmail
	InvalidEmailOrPassword = errors.New("邮箱或密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	//查询用户
	ur, err := svc.repo.FindByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.User{}, InvalidEmailOrPassword
		}
		return domain.User{}, err
	}

	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(ur.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, InvalidEmailOrPassword
	}
	return ur, nil
}
