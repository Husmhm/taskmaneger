package validator

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
	"taskmaneger/param"
	"taskmaneger/pkg/errmsg"
)

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`))),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(func(value interface{}) error {
				return v.checkPhoneNumberUniqueness(value)
			}),
		))

	return err
}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)
	// check uniqueness of phone number
	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {

		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)

		}

	}
	return nil
}

func (v Validator) ValidateLoginRequest(req param.LoginRequest) error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.doesPhoneNumberExist)),
		validation.Field(&req.Password, validation.Required),
	)

	return err
}

func (v Validator) doesPhoneNumberExist(value interface{}) error {
	phoneNumber := value.(string)
	// check uniqueness of phone number
	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf(errmsg.ErrorMsgNotFound)
	}
	return nil
}
