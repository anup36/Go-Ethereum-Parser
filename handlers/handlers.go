package handlers

import (
	"encoding/json"
	"net/http"

	"eth-parser/parser"
)

func SubscribeHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address parameter is required", http.StatusBadRequest)
			return
		}

		success := parser.Subscribe(address)
		if success {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Subscription successful"))
		} else {
			http.Error(w, "Already subscribed", http.StatusConflict)
		}
	}
}

func TransactionsHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "address parameter is required", http.StatusBadRequest)
			return
		}

		transactions := parser.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	}
}

func CurrentBlockHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		currentBlock := parser.GetCurrentBlock()
		json.NewEncoder(w).Encode(currentBlock)
	}
}
