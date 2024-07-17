package main

import (
	"log"
	"net/http"

	"eth-parser/handlers"
	"eth-parser/parser"
	"eth-parser/scheduler"
)

func main() {
	parser := parser.NewEthParser()
	go scheduler.StartScheduler(parser)

	http.HandleFunc("/subscribe", handlers.SubscribeHandler(parser))
	http.HandleFunc("/transactions", handlers.TransactionsHandler(parser))
	http.HandleFunc("/current-block", handlers.CurrentBlockHandler(parser))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
