package ui

import (
	"fmt"
	"strconv"

	"finance_go/models"
	"finance_go/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// MainWindow represents the main application window
type MainWindow struct {
	window              fyne.Window
	financeService      *services.FinanceService
	importExportService *services.ImportExportService
	pdfExportService    *services.PDFExportService
	balance             binding.Float
	transactions        *widget.Table
	amountEntry         *widget.Entry
	descriptionEntry    *widget.Entry
	categoryEntry       *widget.Entry
	typeSelect          *widget.Select
}

// NewMainWindow creates a new main window
func NewMainWindow(window fyne.Window, financeService *services.FinanceService) *MainWindow {
	mw := &MainWindow{
		window:              window,
		financeService:      financeService,
		importExportService: services.NewImportExportService(financeService),
		pdfExportService:    services.NewPDFExportService(financeService),
		balance:             binding.NewFloat(),
	}

	mw.buildUI()
	return mw
}

// buildUI creates the user interface
func (mw *MainWindow) buildUI() {
	// Create form fields
	mw.amountEntry = widget.NewEntry()
	mw.amountEntry.SetPlaceHolder("Valor")

	mw.descriptionEntry = widget.NewEntry()
	mw.descriptionEntry.SetPlaceHolder("Descrição")

	mw.categoryEntry = widget.NewEntry()
	mw.categoryEntry.SetPlaceHolder("Categoria")

	mw.typeSelect = widget.NewSelect([]string{"Receita", "Despesa"}, func(string) {})
	mw.typeSelect.SetSelected("Receita")

	balanceLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(mw.balance, "Saldo: R$ %.2f"))

	// Create transactions table
	mw.transactions = widget.NewTable(
		func() (int, int) {
			return len(mw.financeService.GetTransactions()), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			transactions := mw.financeService.GetTransactions()
			if id.Row < len(transactions) {
				tx := transactions[id.Row]
				switch id.Col {
				case 0:
					label.SetText(tx.Date.Format("02/01/2006"))
				case 1:
					label.SetText(tx.Type)
				case 2:
					label.SetText(fmt.Sprintf("R$ %.2f", tx.Value))
				case 3:
					label.SetText(tx.Description)
				case 4:
					label.SetText(tx.Category)
				}
			}
		},
	)

	// Set table column widths
	mw.transactions.SetColumnWidth(0, 100) // Date
	mw.transactions.SetColumnWidth(1, 80)  // Type
	mw.transactions.SetColumnWidth(2, 100) // Value
	mw.transactions.SetColumnWidth(3, 200) // Description
	mw.transactions.SetColumnWidth(4, 120) // Category

	// Create table headers
	headers := container.NewGridWithColumns(5,
		widget.NewLabel("Data"),
		widget.NewLabel("Tipo"),
		widget.NewLabel("Valor"),
		widget.NewLabel("Descrição"),
		widget.NewLabel("Categoria"),
	)

	// Create table container with headers
	tableContainer := container.NewVBox(
		headers,
		widget.NewSeparator(),
		container.NewVScroll(container.NewMax(mw.transactions)), // Scrollable container that expands
	)

	// Set minimum size for table container to increase height
	tableContainer.Resize(fyne.NewSize(600, 800))

	// Create buttons
	addButton := widget.NewButton("Adicionar", mw.addTransaction)
	importCSVButton := widget.NewButton("Importar CSV", mw.importCSV)
	importExcelButton := widget.NewButton("Importar Excel", mw.importExcel)
	exportCSVButton := widget.NewButton("Exportar CSV", mw.exportCSV)
	exportExcelButton := widget.NewButton("Exportar Excel", mw.exportExcel)
	exportPDFButton := widget.NewButton("Exportar PDF", mw.exportPDF)

	// Create import/export buttons layout - place them at the top
	importExportButtons := container.NewHBox(
		importCSVButton,
		importExcelButton,
		exportCSVButton,
		exportExcelButton,
		exportPDFButton,
	)

	// Create form layout with import/export buttons at the top
	formFields := container.NewGridWithColumns(4,
		mw.amountEntry,
		mw.descriptionEntry,
		mw.categoryEntry,
		mw.typeSelect,
	)

	form := container.NewVBox(
		importExportButtons,
		widget.NewSeparator(),
		formFields,
		addButton,
		balanceLabel,
	)

	// Create main layout using Border layout to make table cover entire remaining size
	// Top: form with import/export buttons
	// Center: table (will expand to fill remaining space)
	// Bottom: small spacer to prevent cutoff
	mainLayout := container.NewBorder(
		container.NewPadded(form),              // top with padding
		container.NewHBox(widget.NewLabel("")), // bottom spacer
		nil,                                    // left
		nil,                                    // right
		container.NewPadded(tableContainer),    // center with padding
	)

	mw.window.SetContent(mainLayout)

	// Update initial balance
	mw.updateBalance()
}

// addTransaction handles adding a new transaction
func (mw *MainWindow) addTransaction() {
	valStr := mw.amountEntry.Text
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		dialog.ShowError(fmt.Errorf("valor inválido"), mw.window)
		return
	}

	typ := mw.typeSelect.Selected
	description := mw.descriptionEntry.Text
	category := mw.categoryEntry.Text

	// Create transaction with new fields
	transaction := models.NewTransaction(typ, val, description, category)
	mw.financeService.AddTransactionFromModel(transaction)

	mw.updateBalance()
	mw.clearForm()
	mw.transactions.Refresh()
}

// clearForm clears all form fields
func (mw *MainWindow) clearForm() {
	mw.amountEntry.SetText("")
	mw.descriptionEntry.SetText("")
	mw.categoryEntry.SetText("")
	mw.typeSelect.SetSelected("Receita")
}

// importCSV handles CSV import
func (mw *MainWindow) importCSV() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		err = mw.importExportService.ImportFromCSV(reader.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("erro ao importar CSV: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Sucesso", "CSV importado com sucesso!", mw.window)
			mw.updateBalance()
			mw.transactions.Refresh()
		}
	}, mw.window)
}

// importExcel handles Excel import
func (mw *MainWindow) importExcel() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		err = mw.importExportService.ImportFromExcel(reader.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("erro ao importar Excel: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Sucesso", "Excel importado com sucesso!", mw.window)
			mw.updateBalance()
			mw.transactions.Refresh()
		}
	}, mw.window)
}

// exportCSV handles CSV export
func (mw *MainWindow) exportCSV() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		err = mw.importExportService.ExportToCSV(writer.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("erro ao exportar CSV: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Sucesso", "CSV exportado com sucesso!", mw.window)
		}
	}, mw.window)
}

// exportExcel handles Excel export
func (mw *MainWindow) exportExcel() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		err = mw.importExportService.ExportToExcel(writer.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("erro ao exportar Excel: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Sucesso", "Excel exportado com sucesso!", mw.window)
		}
	}, mw.window)
}

// exportPDF handles PDF export
func (mw *MainWindow) exportPDF() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		err = mw.pdfExportService.ExportToPDF(writer.URI().Path())
		if err != nil {
			dialog.ShowError(fmt.Errorf("erro ao exportar PDF: %v", err), mw.window)
		} else {
			dialog.ShowInformation("Sucesso", "PDF exportado com sucesso!", mw.window)
		}
	}, mw.window)
}

// updateBalance updates the displayed balance
func (mw *MainWindow) updateBalance() {
	balance := mw.financeService.GetBalance()
	mw.balance.Set(balance)
}

// Refresh refreshes the UI components
func (mw *MainWindow) Refresh() {
	mw.updateBalance()
	mw.transactions.Refresh()
}
