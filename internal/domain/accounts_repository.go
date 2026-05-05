package domain

type IAccountRepository interface {
	Create(account *Account) error
}
