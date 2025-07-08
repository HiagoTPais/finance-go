# Finance Go

A simple financial management application built with Go and Fyne GUI framework.

## Architecture

The application follows a clean layered architecture:

### Layers

- models/: Data structures and basic data operations
  - `transaction.go`: Transaction struct and TransactionList with basic operations

- services/: Business logic layer
  - `finance_service.go`: Handles financial calculations, transaction processing, and business rules
  - `import_export_service.go`: Handles CSV and Excel import/export operations
  - `pdf_export_service.go`: Handles PDF report generation

- ui/: User interface layer using Fyne
  - `main_window.go`: Main application window and UI components

- storage/: Data persistence layer
  - `json_storage.go`: JSON file-based storage implementation

### Data Flow

1. UI Layer → Service Layer: UI components call service methods for business operations
2. Service Layer → Model Layer: Services use models for data manipulation
3. Service Layer → Storage Layer: Services interact with storage for persistence
4. Storage Layer → Model Layer: Storage loads/saves model data

## Features

- Add income and expense transactions with description and category
- Real-time balance calculation
- Transaction history display with date, type, value, description, and category
- Automatic data persistence (JSON)
- Import/Export functionality:
  - Import from CSV files (bank statements)
  - Import from Excel files
  - Export to CSV format
  - Export to Excel format
  - Export to PDF reports with summaries and category breakdowns
- Clean separation of concerns

## Usage

### Running the Application

```bash
go run main.go
```

### Controls

- Add Transaction: Enter amount, description, category, and select type (Receita/Despesa), then click "Adicionar"
- Import Data: Use "Importar CSV" or "Importar Excel" buttons to import bank statements
- Export Data: Use export buttons to save data in various formats
- Save Data: Press ESC key to manually save data
- Auto-save: Data is automatically saved when the application closes

### Import Format

The application expects CSV/Excel files with the following columns:
- Data: Date in YYYY-MM-DD format
- Valor: Amount (positive for income, negative for expenses)
- Descrição: Transaction description
- Categoria: Transaction category

Example CSV format:
```csv
Data,Valor,Descrição,Categoria
2024-01-15,1500.00,Salário,Receita
2024-01-16,-250.00,Supermercado,Alimentação
```

### Export Features

- CSV Export: Exports all transactions in CSV format
- Excel Export: Exports all transactions in Excel format with formatted columns
- PDF Export: Generates comprehensive reports including:
  - Summary with total balance, income, and expenses
  - Complete transaction list
  - Category breakdown with income/expense totals per category

### Data Storage

Transaction data is stored in `data/transactions.json` in JSON format.

## Instalation

```bash
git clone https://github.com/HiagoTPais/finance_go.git
cd finance_go
go mod tidy
go run .
