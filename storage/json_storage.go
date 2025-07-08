package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"finance_go/models"
)

// JSONStorage handles data persistence using JSON files
type JSONStorage struct {
	filePath string
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(filename string) *JSONStorage {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
	}

	return &JSONStorage{
		filePath: filepath.Join(dataDir, filename),
	}
}

// Save saves transaction list to JSON file
func (js *JSONStorage) Save(transactionList *models.TransactionList) error {
	data, err := json.MarshalIndent(transactionList, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling transactions: %w", err)
	}

	err = os.WriteFile(js.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// Load loads transaction list from JSON file
func (js *JSONStorage) Load() (*models.TransactionList, error) {
	data, err := os.ReadFile(js.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty transaction list if file doesn't exist
			return &models.TransactionList{
				Transactions: make([]models.Transaction, 0),
			}, nil
		}
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var transactionList models.TransactionList
	err = json.Unmarshal(data, &transactionList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transactions: %w", err)
	}

	return &transactionList, nil
}
