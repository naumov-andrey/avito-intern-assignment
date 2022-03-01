package service

type Service struct {
	Balance     *BalanceService
	Transaction *TransactionService
}

func NewService(balance *BalanceService, transaction *TransactionService) *Service {
	return &Service{balance, transaction}
}
