package service

import (
	"errors"
	"fmt"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

type AccountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	accessKey       string
}

func NewAccountService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	accessKey string,
) *AccountService {
	return &AccountService{accountRepo, transactionRepo, accessKey}
}

func (s *AccountService) WithTx(tx *gorm.DB) *AccountService {
	newService := *s
	newService.accountRepo = newService.accountRepo.WithTx(tx)
	newService.transactionRepo = newService.transactionRepo.WithTx(tx)
	return &newService
}

func (s *AccountService) GetBalance(userId int) (decimal.Decimal, error) {
	return s.accountRepo.GetBalance(userId)
}

func (s *AccountService) GetConvertedBalance(userId int, currency string) (decimal.Decimal, error) {
	balance, err := s.accountRepo.GetBalance(userId)
	if err != nil {
		return balance, err
	}

	apiCall := fmt.Sprintf("https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=RUB", s.accessKey)
	response, err := http.Get(apiCall)
	if err != nil {
		return balance, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return balance, err
	}

	rateStr := gjson.Get(string(responseData), "data."+currency).String()
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return balance, errors.New(fmt.Sprintf("Error while getting RUB/%s currecny rate", currency))
	}

	return balance.Mul(decimal.NewFromFloat(rate)), nil
}

func (s *AccountService) UpdateBalance(userId int, amount decimal.Decimal, description string) (decimal.Decimal, error) {
	if amount.Equal(decimal.Decimal{}) {
		return s.accountRepo.GetBalance(userId)
	}

	var zero decimal.Decimal
	balance, err := s.accountRepo.GetBalance(userId)
	if err != nil {
		return zero, err
	}

	newBalance := balance.Add(amount)
	if newBalance.LessThan(zero) {
		return zero, errors.New("account balance is less than credit amount")
	}

	err = s.accountRepo.UpdateBalance(userId, newBalance)
	if err != nil {
		return zero, err
	}

	if _, err = s.transactionRepo.CreateTransaction(userId, amount, description); err != nil {
		return zero, err
	}

	balance, err = s.accountRepo.GetBalance(userId)
	if err != nil {
		return zero, err
	}

	return balance, nil
}

func (s *AccountService) GetHistory(
	userId int,
	limit int,
	cursor int,
	sortBy string,
	orderBy string,
) (model.HistoryOutput, error) {
	if sortBy != "date" && sortBy != "amount" {
		return model.HistoryOutput{}, errors.New("sort must be by date or amount")
	}

	if orderBy != "asc" && orderBy != "desc" {
		return model.HistoryOutput{}, errors.New("order must be asc or desc")
	}

	var (
		data []model.Transaction
		err  error
	)
	if cursor != -1 {
		data, err = s.transactionRepo.GetHistoryWithCursor(userId, limit, cursor, sortBy, orderBy)
	} else {
		data, err = s.transactionRepo.GetHistory(userId, limit, sortBy, orderBy)
	}

	if err != nil {
		return model.HistoryOutput{}, err
	}

	if len(data) != 0 && limit == len(data) {
		cursor = data[len(data)-1].Id
	} else {
		cursor = -1
	}

	return model.HistoryOutput{Data: data, Cursor: cursor}, nil
}

func (s *AccountService) CreateTransfer(transferData model.TransferData) (model.TransferOutput, error) {
	creditUserBalance, err := s.accountRepo.GetBalance(transferData.CreditUserId)
	if err != nil {
		return model.TransferOutput{}, err
	}

	if creditUserBalance.LessThan(transferData.Amount) {
		return model.TransferOutput{}, errors.New("credit account balance is less than amount")
	}

	credit, err := s.transactionRepo.CreateTransaction(
		transferData.CreditUserId,
		transferData.Amount.Neg(),
		transferData.Description,
	)
	if err != nil {
		return model.TransferOutput{}, err
	}

	newCreditBalance := creditUserBalance.Sub(transferData.Amount)
	if err = s.accountRepo.UpdateBalance(transferData.CreditUserId, newCreditBalance); err != nil {
		return model.TransferOutput{}, err
	}

	debitUserBalance, err := s.accountRepo.GetBalance(transferData.DebitUserId)
	if err != nil {
		return model.TransferOutput{}, err
	}

	debit, err := s.transactionRepo.CreateTransaction(
		transferData.DebitUserId,
		transferData.Amount,
		transferData.Description,
	)
	if err != nil {
		return model.TransferOutput{}, err
	}

	newDebitBalance := debitUserBalance.Add(transferData.Amount)
	if err = s.accountRepo.UpdateBalance(transferData.DebitUserId, newDebitBalance); err != nil {
		return model.TransferOutput{}, err
	}

	return model.TransferOutput{Debit: debit, Credit: credit}, nil
}
