package apperr

import "errors"

var (
	AccountAlreadyExists   = errors.New("account already exists")
	AccountNotFound        = errors.New("account not found")
	OrderValidationFailed  = errors.New("order validation error")
	OrderAlreadyRegistered = errors.New("order already registered")
	InsufficientBalance    = errors.New("insufficient balance")
	AccountInvalidBalance  = errors.New("account invalid balance")
	InvalidCreditValue     = errors.New("invalid credit value")
	InvalidDebitValue      = errors.New("invalid debit value")
	AccrualNotFoundError   = errors.New("accrual not found")
	InternalServerError    = errors.New("internal server error")
	TooManyRequestsError   = errors.New("too many requests")
)
