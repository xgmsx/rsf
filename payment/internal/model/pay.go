package model

type PaymentMethod int32

const (
	PaymentMethod_UNSPECIFIED    PaymentMethod = 0
	PaymentMethod_CARD           PaymentMethod = 1
	PaymentMethod_SBP            PaymentMethod = 2
	PaymentMethod_CREDIT_CARD    PaymentMethod = 3
	PaymentMethod_INVESTOR_MONEY PaymentMethod = 4
)
