package main

import (
	"fmt"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var amountEntry *gtk.Entry
var typeCombo *gtk.ComboBoxText
var balanceLabel *gtk.Label
var listStore *gtk.ListStore

func buildMainWindow() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Financeiro")
	win.SetDefaultSize(400, 400)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	win.Add(vbox)

	amountEntry, _ = gtk.EntryNew()
	amountEntry.SetPlaceholderText("Valor")
	vbox.PackStart(amountEntry, false, false, 0)

	typeCombo, _ = gtk.ComboBoxTextNew()
	typeCombo.AppendText("Receita")
	typeCombo.AppendText("Despesa")
	typeCombo.SetActive(0)
	vbox.PackStart(typeCombo, false, false, 0)

	addButton, _ := gtk.ButtonNewWithLabel("Adicionar")
	addButton.Connect("clicked", onAddTransaction)
	vbox.PackStart(addButton, false, false, 0)

	balanceLabel, _ = gtk.LabelNew("Saldo: R$ 0.00")
	vbox.PackStart(balanceLabel, false, false, 0)

	listStore, _ = gtk.ListStoreNew(
		glib.TYPE_STRING,
		glib.TYPE_DOUBLE,
	)

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()
	column, _ := gtk.TreeViewColumnNewWithAttribute("Tipo", renderer, "text", 0)
	treeView.AppendColumn(column)

	renderer2, _ := gtk.CellRendererTextNew()
	column2, _ := gtk.TreeViewColumnNewWithAttribute("Valor", renderer2, "text", 1)
	treeView.AppendColumn(column2)

	vbox.PackStart(treeView, true, true, 0)

	win.ShowAll()
}

func onAddTransaction() {
	text, _ := amountEntry.GetText()
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		fmt.Println("Erro ao converter valor")
		return
	}

	typ := typeCombo.GetActiveText()
	if typ == "Despesa" {
		value = -value
	}

	addTransactionData(value)

	iter := listStore.Append()
	listStore.Set(iter, []int{0, 1}, []interface{}{typ, value})

	saldo := getBalance()
	balanceLabel.SetText(fmt.Sprintf("Saldo: R$ %.2f", saldo))

	amountEntry.SetText("")
}
