package repository

import (
	"context"
	"jike/internal/domain"
	"jike/internal/repository/cache"
	"jike/internal/repository/dao"
)

var (
	ErrDuplicateEmail      = dao.ErrDuplicateEmail
	InvalidEmailOrPassword = dao.InvalidEmailOrPassword
	ErrUserNotFound        = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cacheObj *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cacheObj,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	switch err {
	case nil:
		return u, nil
	case cache.ErrKeyNotExist:
		ue, err := r.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		u = domain.User{
			Id:       ue.Id,
			Email:    ue.Email,
			Password: ue.Password,
		}
		_ = r.cache.Set(ctx, u) //记录日志
		return u, err
	default:
		return domain.User{}, err
	}
}
