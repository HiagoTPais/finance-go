package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"finance_go/models"

	"github.com/xuri/excelize/v2"
)

// ImportExportService handles data import and export operations
type ImportExportService struct {
	financeService *FinanceService
}

// NewImportExportService creates a new import/export service
func NewImportExportService(financeService *FinanceService) *ImportExportService {
	return &ImportExportService{
		financeService: financeService,
	}
}

// ImportFromCSV imports transactions from a CSV file
func (ies *ImportExportService) ImportFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header and one data row")
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 4 {
			continue // Skip incomplete rows
		}

		// Parse date
		date, err := time.Parse("2006-01-02", strings.TrimSpace(record[0]))
		if err != nil {
			date = time.Now() // Use current date if parsing fails
		}

		// Parse amount
		amountStr := strings.ReplaceAll(strings.TrimSpace(record[1]), ",", ".")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue // Skip rows with invalid amounts
		}

		// Determine transaction type based on amount
		transactionType := "Receita"
		if amount < 0 {
			transactionType = "Despesa"
			amount = -amount // Make positive for internal storage
		}

		description := strings.TrimSpace(record[2])
		category := strings.TrimSpace(record[3])

		transaction := models.NewTransactionWithDate(transactionType, amount, description, category, date)
		ies.financeService.AddTransactionFromModel(transaction)
	}

	return nil
}

// ImportFromExcel imports transactions from an Excel file
func (ies *ImportExportService) ImportFromExcel(filename string) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("error opening Excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error reading Excel sheet: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("Excel file must have at least a header and one data row")
	}

	// Skip header row
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 4 {
			continue // Skip incomplete rows
		}

		// Parse date
		dateStr := strings.TrimSpace(row[0])
		var date time.Time
		if dateStr != "" {
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				date = time.Now() // Use current date if parsing fails
			}
		} else {
			date = time.Now()
		}

		// Parse amount
		amountStr := strings.ReplaceAll(strings.TrimSpace(row[1]), ",", ".")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue // Skip rows with invalid amounts
		}

		// Determine transaction type based on amount
		transactionType := "Receita"
		if amount < 0 {
			transactionType = "Despesa"
			amount = -amount // Make positive for internal storage
		}

		description := strings.TrimSpace(row[2])
		category := strings.TrimSpace(row[3])

		transaction := models.NewTransactionWithDate(transactionType, amount, description, category, date)
		ies.financeService.AddTransactionFromModel(transaction)
	}

	return nil
}

// ExportToCSV exports transactions to a CSV file
func (ies *ImportExportService) ExportToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Data", "Valor", "Descrição", "Categoria", "Tipo"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data
	transactions := ies.financeService.GetTransactions()
	for _, tx := range transactions {
		amount := tx.Value
		if tx.Type == "Despesa" {
			amount = -amount
		}

		record := []string{
			tx.Date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", amount),
			tx.Description,
			tx.Category,
			tx.Type,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	return nil
}

// ExportToExcel exports transactions to an Excel file
func (ies *ImportExportService) ExportToExcel(filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	// Set headers
	headers := []string{"Data", "Valor", "Descrição", "Categoria", "Tipo"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Write data
	transactions := ies.financeService.GetTransactions()
	for i, tx := range transactions {
		row := i + 2 // Start from row 2 (after header)
		amount := tx.Value
		if tx.Type == "Despesa" {
			amount = -amount
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), tx.Date.Format("2006-01-02"))
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), amount)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), tx.Description)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), tx.Category)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), tx.Type)
	}

	// Auto-fit columns
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth("Sheet1", col, col, 15)
	}

	return f.SaveAs(filename)
}
