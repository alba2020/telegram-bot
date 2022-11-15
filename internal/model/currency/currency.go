package currency

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

const (
	USD = "USD"
	CNY = "CNY"
	EUR = "EUR"
	RUB = "RUB"
)

type CurrencyRates struct {
	Base  string `json:"base"`
	Rates struct {
		USD float64 `json:"USD"`
		EUR float64 `json:"EUR"`
		CNY float64 `json:"CNY"`
	} `json:"rates"`
}

type Service struct {
	valid []string
	rates map[string]float64 // rub
}

const CBR = "https://www.cbr-xml-daily.ru/latest.js"

func New() (*Service, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(CBR)
	if err != nil {
		logger.Error("No response from currency site")
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body) // []byte
	defer resp.Body.Close()

	currencyRates := CurrencyRates{}
	if err := json.Unmarshal(body, &currencyRates); err != nil {
		logger.Error("Can not unmarshal JSON")
		return nil, err
	}

	svc := Service{
		valid: []string{USD, CNY, EUR, RUB},
		rates: map[string]float64{
			USD: currencyRates.Rates.USD,
			EUR: currencyRates.Rates.EUR,
			CNY: currencyRates.Rates.CNY,
			RUB: 1.0000,
		},
	}

	logger.Info("Currency service ready")
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
