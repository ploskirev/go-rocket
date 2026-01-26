package model

type PaymentInfo struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

type PaymentMethod int32

const (
	PaymentMethod_UNKNOWN PaymentMethod = iota
	PaymentMethod_CARD
	PaymentMethod_SBP
	PaymentMethod_CREDIT_CARD
	PaymentMethod_INVESTOR_MONEY
)
