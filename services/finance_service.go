package services

import "finance_go/models"

// FinanceService handles business logic for financial operations
type FinanceService struct {
	transactionList *models.TransactionList
}

// NewFinanceService creates a new finance service
func NewFinanceService() *FinanceService {
	return &FinanceService{
		transactionList: &models.TransactionList{
			Transactions: make([]models.Transaction, 0),
		},
	}
}

// AddTransaction adds a new transaction with business logic
func (fs *FinanceService) AddTransaction(transactionType string, amount float64) {
	if transactionType == "Despesa" {
		amount = -amount
	}

	transaction := models.NewTransaction(transactionType, amount, "", "")
	fs.transactionList.AddTransaction(transaction)
}

// AddTransactionFromModel adds a transaction directly from a model (for imports)
func (fs *FinanceService) AddTransactionFromModel(transaction models.Transaction) {
	fs.transactionList.AddTransaction(transaction)
}

// GetTransactions returns all transactions
func (fs *FinanceService) GetTransactions() []models.Transaction {
	return fs.transactionList.GetTransactions()
}

// GetBalance returns the current balance
func (fs *FinanceService) GetBalance() float64 {
	return fs.transactionList.GetBalance()
}

// GetTransactionList returns the transaction list for storage operations
func (fs *FinanceService) GetTransactionList() *models.TransactionList {
	return fs.transactionList
}

// SetTransactionList sets the transaction list (for loading from storage)
func (fs *FinanceService) SetTransactionList(tl *models.TransactionList) {
	fs.transactionList = tl
}
