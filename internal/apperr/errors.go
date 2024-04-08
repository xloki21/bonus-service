package apperr

import "errors"

var (
	AccountAlreadyExists              = errors.New("account already exists")
	AccountNotFound                   = errors.New("account not found")
	OrderValidationFailed             = errors.New("order validation error")
	OrderAlreadyRegistered            = errors.New("order already registered")
	InsufficientBalance               = errors.New("insufficient balance")
	AccountInvalidBalance             = errors.New("account invalid balance")
	AccrualServiceTooManyRequests     = errors.New("too many requests to accrual service")
	AccrualNotFound                   = errors.New("accrual not found")
	AccrualServiceInternalServerError = errors.New("accrual service internal server error")
)
