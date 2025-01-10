package ArifpayGoPlugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Payment struct {
	APIKey     string
	ExpireDate time.Time
}

type PaymentRequest struct {
	CancelUrl      string        `json:"cancelUrl"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email"`
	Nonce          string        `json:"nonce"`
	ErrorUrl       string        `json:"errorUrl"`
	NotifyUrl      string        `json:"notifyUrl"`
	SuccessUrl     string        `json:"successUrl"`
	PaymentMethods []string      `json:"paymentMethods"`
	ExpireDate     time.Time     `json:"expireDate"`
	Items          []interface{} `json:"items"`
	Beneficiaries  []struct {
		AccountNumber string  `json:"accountNumber"`
		Bank          string  `json:"bank"`
		Amount        float64 `json:"amount"`
	} `json:"beneficiaries"`
	Lang string `json:"lang"`
}

func NewPayment(apiKey string, expireDate time.Time) *Payment {
	return &Payment{
		APIKey:     apiKey,
		ExpireDate: expireDate,
	}
}

func (p *Payment) MakePayment(paymentRequest PaymentRequest) (string, error) {
	paymentRequestBytes, err := json.Marshal(paymentRequest)
	if err != nil {
		return "", err
	}
	paymentRequest.Nonce = fmt.Sprintf("%d", time.Now().UnixNano())
	paymentRequest.ExpireDate = p.ExpireDate

	req, err := http.NewRequest("POST", "http://196.189.44.37:2000/api/sandbox/checkout/session", bytes.NewBuffer(paymentRequestBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-arifpay-key", p.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}
