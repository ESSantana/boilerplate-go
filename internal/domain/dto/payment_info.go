package dto

import "strings"

type Item struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  int     `json:"quantity"`
}

type PaymentInfo struct {
	CustomerName      string  `json:"customer_name"`
	CustomerEmail     string  `json:"customer_email"`
	CustomerPhone     string  `json:"customer_phone"`
	Items             []Item  `json:"items"`
}

func (pi *PaymentInfo) GetFirstName() string {
	return strings.ToUpper(strings.Split(pi.CustomerName, " ")[0])
}

func (pi *PaymentInfo) GetLastName() string {
	nameSlice := strings.Split(pi.CustomerName, " ")
	return strings.ToUpper(nameSlice[len(nameSlice)-1])
}

func (pi *PaymentInfo) GetPhoneAreaCode() string {
	return pi.CustomerPhone[:2]
}

func (pi *PaymentInfo) GetPhoneNumber() string {
	return pi.CustomerPhone[2:]
}
