package parser

import "eth-parser/clients"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []clients.Transaction
}
