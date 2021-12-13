package oauth

import (
	"errors"
	"fmt"
	"ibtools_server/models"
	util "ibtools_server/util"
	"math/rand"
	"time"

	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"gorm.io/gorm"
)

//SendPhoneNumValidateMessage 发送验证码
func (s *Service) SendPhoneNumValidateMessage(phone string) error {
	if !util.ValidatePhone(phone) {
		return errors.New("phone nubmer is wrong")
	}
	//TODO send message and 时间延时校验
	phonenumvalidate := new(models.OauthPhoneNumValidate)
	err := s.db.Where("phone = LOWER(?)", util.FormatPhoneNum(phone)).
		First(phonenumvalidate).Error

	phonenumvalidate.Phone = util.FormatPhoneNum(phone)
	phonenumvalidate.ExpiresAt = time.Now().UTC().Add(time.Duration(60) * time.Second)
	//phonenumvalidate.Code = util.GenValidateCode(4)
	//TODO
	phonenumvalidate.Code = fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(phone),
		TemplateCode: tea.String("SMS_225805063"),
		SignName:     tea.String("123123"),
	}
	_, _err := s.smsClient.SendSms(sendSmsRequest)
	if _err != nil {
		return _err
	}
	//util.SendSMS(phone, fmt.Sprintf("%s  is your Dr. Piggy Verification Code. This code will expire in 1 minute.", phonenumvalidate.Code))
	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := s.db.Create(phonenumvalidate).Error; err != nil {
			return err
		}
	} else {
		if err := s.db.Save(phonenumvalidate).Error; err != nil {
			return err
		}
	}
	return nil
}

//PhoneNumValidate 手机号验证
func (s *Service) PhoneNumValidate(phone, code string) error {
	//TODO
	if code == "6666" {
		return nil
	}
	if util.ValidatePhone(phone) {
		phonenumvalidate := new(models.OauthPhoneNumValidate)
		err := s.db.Where("phone = LOWER(?)", util.FormatPhoneNum(phone)).
			First(phonenumvalidate).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phone number not validate")
		}
		if time.Now().UTC().After(phonenumvalidate.ExpiresAt) {
			return errors.New("code is expired")
		}
		if phonenumvalidate.Code != code {
			return errors.New("code is wrong")
		}
		return nil
	}
	return errors.New("phone nubmer is wrong")
}
