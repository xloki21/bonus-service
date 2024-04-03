package accrual

import (
	"context"
	"errors"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"github.com/xloki21/bonus-service/internal/entity/transaction"
	"github.com/xloki21/bonus-service/internal/repository"
)

type Accrual interface {
	RequestOrderReward(ctx context.Context, order *order.Order) (uint, error)
}

type Service struct {
	repo repository.Transaction
}

// RequestOrderReward requests order reward.
func (a *Service) RequestOrderReward(ctx context.Context, order *order.Order) (uint, error) {

	transactions, err := a.repo.GetOrderTransactions(ctx, order)
	if err != nil {
		return 0, errors.New("accrual transactions not found")
	}

	var reward uint = 0
	for _, tx := range transactions {
		if tx.Status == transaction.PROCESSED {
			reward += tx.Reward
		} else {
			return 0, errors.New("accrual is not completed yet")
		}
	}

	return reward, nil
}

func NewAccrualService(repo repository.Transaction) *Service {
	return &Service{repo: repo}
}
