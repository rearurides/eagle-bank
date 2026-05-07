package domain

type ITransactionsRepository interface {
	Deposit(tan *Transaction) error
	Withdraw(tan *Transaction) error
}
