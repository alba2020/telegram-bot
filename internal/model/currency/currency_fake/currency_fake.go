package currencyfake

import (
	"errors"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

const (
	USD = "USD"
	CNY = "CNY"
	EUR = "EUR"
	RUB = "RUB"
)

type Service struct {
	valid []string
	rates map[string]float64 // rub
}

func New() (*Service, error) {
	svc := Service{
		valid: []string{USD, CNY, EUR, RUB},
		rates: map[string]float64{
			USD: 100.00,
			EUR: 105.00,
			CNY: 10.00,
			RUB: 1.00,
		},
	}

	logger.Info("Currency service fake ready")
	return &svc, nil
}

func (svc *Service) IsValid(val string) bool {
	for _, cur := range svc.valid {
		if val == cur {
			return true
		}
	}
	return false
}

func (svc *Service) GetValid() []string {
	return svc.valid
}

func (svc *Service) Convert(
	amount float64, from string, to string) (float64, error) {

	if !svc.IsValid(from) || !svc.IsValid(to) {
		return 0, errors.New("bad currency code")
	}
	rub := amount / svc.rates[from]
	res := rub * svc.rates[to]

	return res, nil
}
