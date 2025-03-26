package interfaces

type SSOManager interface { 
	GetProvider(provider string) (SSOProvider, error)
}