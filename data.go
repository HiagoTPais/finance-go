package main

var balance float64 = 0

func addTransactionData(amount float64) {
	balance += amount
}

func getBalance() float64 {
	return balance
}
