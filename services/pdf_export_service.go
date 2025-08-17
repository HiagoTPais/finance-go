package services

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PDFExportService handles PDF report generation
type PDFExportService struct {
	financeService *FinanceService
}

// NewPDFExportService creates a new PDF export service
func NewPDFExportService(financeService *FinanceService) *PDFExportService {
	return &PDFExportService{
		financeService: financeService,
	}
}

// ExportToPDF exports transactions to a PDF report
func (pes *PDFExportService) ExportToPDF(filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(190, 10, "Relatório Financeiro")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Resumo")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	balance := pes.financeService.GetBalance()
	pdf.Cell(190, 6, fmt.Sprintf("Saldo Total: R$ %.2f", balance))
	pdf.Ln(8)

	transactions := pes.financeService.GetTransactions()
	totalIncome := 0.0
	totalExpenses := 0.0

	for _, tx := range transactions {
		if tx.Type == "Receita" {
			totalIncome += tx.Value
		} else {
			totalExpenses += tx.Value
		}
	}

	pdf.Cell(190, 6, fmt.Sprintf("Total Receitas: R$ %.2f", totalIncome))
	pdf.Ln(6)
	pdf.Cell(190, 6, fmt.Sprintf("Total Despesas: R$ %.2f", totalExpenses))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Transações")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 9)
	headers := []string{"Data", "Tipo", "Valor", "Descrição", "Categoria"}
	widths := []float64{25, 20, 25, 70, 30}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, tx := range transactions {
		amount := tx.Value
		if tx.Type == "Despesa" {
			amount = -amount
		}

		pdf.CellFormat(widths[0], 6, tx.Date.Format("02/01/2006"), "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[1], 6, tx.Type, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], 6, fmt.Sprintf("R$ %.2f", amount), "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[3], 6, tx.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[4], 6, tx.Category, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Resumo por Categoria")
	pdf.Ln(10)

	categories := pes.financeService.GetTransactionList().GetCategories()
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(60, 7, "Categoria")
	pdf.Cell(40, 7, "Receitas")
	pdf.Cell(40, 7, "Despesas")
	pdf.Cell(40, 7, "Saldo")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, category := range categories {
		categoryTxs := pes.financeService.GetTransactionList().GetTransactionsByCategory(category)
		categoryIncome := 0.0
		categoryExpenses := 0.0

		for _, tx := range categoryTxs {
			if tx.Type == "Receita" {
				categoryIncome += tx.Value
			} else {
				categoryExpenses += tx.Value
			}
		}

		categoryBalance := categoryIncome - categoryExpenses

		pdf.Cell(60, 6, category)
		pdf.Cell(40, 6, fmt.Sprintf("R$ %.2f", categoryIncome))
		pdf.Cell(40, 6, fmt.Sprintf("R$ %.2f", categoryExpenses))
		pdf.Cell(40, 6, fmt.Sprintf("R$ %.2f", categoryBalance))
		pdf.Ln(-1)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(190, 6, fmt.Sprintf("Relatório gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))

	return pdf.OutputFileAndClose(filename)
}

// ExportMonthlyReport exports a monthly report to PDF
func (pes *PDFExportService) ExportMonthlyReport(filename string, year int, month time.Month) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(190, 10, fmt.Sprintf("Relatório Mensal - %s/%d", month.String(), year))
	pdf.Ln(15)

	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	monthlyTransactions := pes.financeService.GetTransactionList().GetTransactionsByDateRange(startDate, endDate)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Resumo Mensal")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	monthlyIncome := 0.0
	monthlyExpenses := 0.0

	for _, tx := range monthlyTransactions {
		if tx.Type == "Receita" {
			monthlyIncome += tx.Value
		} else {
			monthlyExpenses += tx.Value
		}
	}

	monthlyBalance := monthlyIncome - monthlyExpenses

	pdf.Cell(190, 6, fmt.Sprintf("Receitas do Mês: R$ %.2f", monthlyIncome))
	pdf.Ln(6)
	pdf.Cell(190, 6, fmt.Sprintf("Despesas do Mês: R$ %.2f", monthlyExpenses))
	pdf.Ln(6)
	pdf.Cell(190, 6, fmt.Sprintf("Saldo do Mês: R$ %.2f", monthlyBalance))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Transações do Mês")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 9)
	headers := []string{"Data", "Tipo", "Valor", "Descrição", "Categoria"}
	widths := []float64{25, 20, 25, 70, 30}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, tx := range monthlyTransactions {
		amount := tx.Value
		if tx.Type == "Despesa" {
			amount = -amount
		}

		pdf.CellFormat(widths[0], 6, tx.Date.Format("02/01/2006"), "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[1], 6, tx.Type, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], 6, fmt.Sprintf("R$ %.2f", amount), "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[3], 6, tx.Description, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[4], 6, tx.Category, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(190, 6, fmt.Sprintf("Relatório gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))

	return pdf.OutputFileAndClose(filename)
}
