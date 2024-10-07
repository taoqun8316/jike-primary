package service

import (
	"context"
	"fmt"
	"jike/internal/repository"
	"jike/internal/service/sms"
	"math/rand"
)

const TplId = "SMS_253000001"

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.sms.Send(ctx, TplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) error {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(100000)
	return fmt.Sprintf("%06d", num)
}
