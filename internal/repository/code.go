package repository

import (
	"context"
	"errors"
	"jike/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cacheObj *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cacheObj,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz string, phone string, verifyCode string) error {
	code, err := repo.cache.Get(ctx, biz, phone)
	if err != nil {
		return err
	}
	if verifyCode == code {
		return nil
	}
	return errors.New("验证码错误")
}
