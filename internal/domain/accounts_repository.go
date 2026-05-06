package domain

type IAccountRepository interface {
	Create(account *Account) error
	GetByAccountNumber(userId, accountNumber string) (*Account, error)
}
