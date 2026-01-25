package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
}

type OrderStatus string

const (
	PENDING_PAYMENT OrderStatus = "PENDING_PAYMENT"
	PAID            OrderStatus = "PAID"
	CANCELLED       OrderStatus = "CANCELLED"
)

type CreateOrderData struct {
	UserUUID  string
	PartUUIDs []string
}

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
