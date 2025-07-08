package models

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// TransactionList holds a collection of transactions
type TransactionList struct {
	Transactions []Transaction `json:"transactions"`
}

// NewTransaction creates a new transaction
func NewTransaction(transactionType string, value float64, description, category string) Transaction {
	return Transaction{
		ID:          generateID(),
		Type:        transactionType,
		Value:       value,
		Description: description,
		Category:    category,
		Date:        time.Now(),
	}
}

// NewTransactionWithDate creates a new transaction with a specific date
func NewTransactionWithDate(transactionType string, value float64, description, category string, date time.Time) Transaction {
	return Transaction{
		ID:          generateID(),
		Type:        transactionType,
		Value:       value,
		Description: description,
		Category:    category,
		Date:        date,
	}
}

// AddTransaction adds a transaction to the list
func (tl *TransactionList) AddTransaction(transaction Transaction) {
	tl.Transactions = append(tl.Transactions, transaction)
}

// GetTransactions returns all transactions
func (tl *TransactionList) GetTransactions() []Transaction {
	return tl.Transactions
}

// GetBalance calculates the total balance from all transactions
func (tl *TransactionList) GetBalance() float64 {
	var total float64
	for _, tx := range tl.Transactions {
		total += tx.Value
	}
	return total
}

// GetTransactionsByDateRange returns transactions within a date range
func (tl *TransactionList) GetTransactionsByDateRange(start, end time.Time) []Transaction {
	var filtered []Transaction
	for _, tx := range tl.Transactions {
		if (tx.Date.Equal(start) || tx.Date.After(start)) && (tx.Date.Equal(end) || tx.Date.Before(end)) {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}

// GetTransactionsByCategory returns transactions filtered by category
func (tl *TransactionList) GetTransactionsByCategory(category string) []Transaction {
	var filtered []Transaction
	for _, tx := range tl.Transactions {
		if tx.Category == category {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}

// GetCategories returns all unique categories
func (tl *TransactionList) GetCategories() []string {
	categories := make(map[string]bool)
	for _, tx := range tl.Transactions {
		categories[tx.Category] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// Simple ID generator (in a real app, you'd use a proper ID system)
var nextID = 1

func generateID() int {
	id := nextID
	nextID++
	return id
}
