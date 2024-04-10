// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	account "github.com/xloki21/bonus-service/internal/entity/account"
	order "github.com/xloki21/bonus-service/internal/entity/order"
	transaction "github.com/xloki21/bonus-service/internal/entity/transaction"
)

// MockOrder is a mock of Order interface.
type MockOrder struct {
	ctrl     *gomock.Controller
	recorder *MockOrderMockRecorder
}

// MockOrderMockRecorder is the mock recorder for MockOrder.
type MockOrderMockRecorder struct {
	mock *MockOrder
}

// NewMockOrder creates a new mock instance.
func NewMockOrder(ctrl *gomock.Controller) *MockOrder {
	mock := &MockOrder{ctrl: ctrl}
	mock.recorder = &MockOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrder) EXPECT() *MockOrderMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockOrder) Register(ctx context.Context, o order.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, o)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockOrderMockRecorder) Register(ctx, o interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockOrder)(nil).Register), ctx, o)
}

// MockAccount is a mock of Account interface.
type MockAccount struct {
	ctrl     *gomock.Controller
	recorder *MockAccountMockRecorder
}

// MockAccountMockRecorder is the mock recorder for MockAccount.
type MockAccountMockRecorder struct {
	mock *MockAccount
}

// NewMockAccount creates a new mock instance.
func NewMockAccount(ctrl *gomock.Controller) *MockAccount {
	mock := &MockAccount{ctrl: ctrl}
	mock.recorder = &MockAccountMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccount) EXPECT() *MockAccountMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAccount) Create(arg0 context.Context, arg1 account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAccountMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAccount)(nil).Create), arg0, arg1)
}

// Credit mocks base method.
func (m *MockAccount) Credit(arg0 context.Context, arg1 account.UserID, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Credit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Credit indicates an expected call of Credit.
func (mr *MockAccountMockRecorder) Credit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Credit", reflect.TypeOf((*MockAccount)(nil).Credit), arg0, arg1, arg2)
}

// Debit mocks base method.
func (m *MockAccount) Debit(arg0 context.Context, arg1 account.UserID, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Debit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Debit indicates an expected call of Debit.
func (mr *MockAccountMockRecorder) Debit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debit", reflect.TypeOf((*MockAccount)(nil).Debit), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockAccount) Delete(arg0 context.Context, arg1 account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccount)(nil).Delete), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockAccount) FindByID(arg0 context.Context, arg1 account.UserID) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockAccountMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAccount)(nil).FindByID), arg0, arg1)
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// FindUnprocessed mocks base method.
func (m *MockTransaction) FindUnprocessed(ctx context.Context, limit int64) ([]transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUnprocessed", ctx, limit)
	ret0, _ := ret[0].([]transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUnprocessed indicates an expected call of FindUnprocessed.
func (mr *MockTransactionMockRecorder) FindUnprocessed(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUnprocessed", reflect.TypeOf((*MockTransaction)(nil).FindUnprocessed), ctx, limit)
}

// GetOrderTransactions mocks base method.
func (m *MockTransaction) GetOrderTransactions(arg0 context.Context, arg1 order.Order) ([]transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderTransactions", arg0, arg1)
	ret0, _ := ret[0].([]transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderTransactions indicates an expected call of GetOrderTransactions.
func (mr *MockTransactionMockRecorder) GetOrderTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderTransactions", reflect.TypeOf((*MockTransaction)(nil).GetOrderTransactions), arg0, arg1)
}

// RewardAccounts mocks base method.
func (m *MockTransaction) RewardAccounts(ctx context.Context, limit int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RewardAccounts", ctx, limit)
	ret0, _ := ret[0].(error)
	return ret0
}

// RewardAccounts indicates an expected call of RewardAccounts.
func (mr *MockTransactionMockRecorder) RewardAccounts(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RewardAccounts", reflect.TypeOf((*MockTransaction)(nil).RewardAccounts), ctx, limit)
}

// Update mocks base method.
func (m *MockTransaction) Update(ctx context.Context, tx *transaction.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTransactionMockRecorder) Update(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTransaction)(nil).Update), ctx, tx)
}
