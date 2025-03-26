package interfaces

type PaymentManager interface {
	NewMercadoPagoProvider() (PaymentProvider, error)
}
