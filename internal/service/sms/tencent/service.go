package tencent

import (
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	SignName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		SignName: &signName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.SignName
	req.TemplateId = &tplId
	req.PhoneNumberSet = stringToPointer(numbers)
	req.TemplateParamSet = stringToPointer(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("短信发送失败: %s, %s ", *(status.Code), *(status.Message))
		}
	}
	return nil
}

func stringToPointer(strs []string) []*string {
	// 创建一个新的 []*string 切片
	ptrStrings := make([]*string, len(strs))

	// 遍历原始切片并创建指针
	for i, s := range strs {
		ptrStrings[i] = &s
	}
	return ptrStrings
}
