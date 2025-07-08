package main

import (
	"log"

	"finance_go/services"
	"finance_go/storage"
	"finance_go/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize the application
	a := app.New()
	w := a.NewWindow("Financeiro")

	// Set window size to 800x800 for more table space
	w.Resize(fyne.NewSize(800, 800))

	// Initialize storage layer
	storage := storage.NewJSONStorage("transactions.json")

	// Load existing data
	transactionList, err := storage.Load()
	if err != nil {
		log.Printf("Error loading data: %v", err)
	}

	// Initialize service layer
	financeService := services.NewFinanceService()
	if transactionList != nil {
		financeService.SetTransactionList(transactionList)
	}

	// Initialize UI layer
	_ = ui.NewMainWindow(w, financeService)

	// Set up auto-save functionality
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			// Save data when ESC is pressed
			err := storage.Save(financeService.GetTransactionList())
			if err != nil {
				log.Printf("Error saving data: %v", err)
			}
		}
	})

	// Show and run the application
	w.ShowAndRun()

	// Save data when application closes
	err = storage.Save(financeService.GetTransactionList())
	if err != nil {
		log.Printf("Error saving data on exit: %v", err)
	}
}
