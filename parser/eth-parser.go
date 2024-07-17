package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"

	"eth-parser/clients"
)

type EthParser struct {
	currentBlock int
	subscribers  map[string][]clients.Transaction
	mu           sync.Mutex
}

func NewEthParser() *EthParser {
	return &EthParser{
		subscribers: make(map[string][]clients.Transaction),
	}
}

func (ep *EthParser) GetCurrentBlock() int {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	return ep.currentBlock
}

func (ep *EthParser) Subscribe(address string) bool {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	address = strings.ToLower(address)
	if _, exists := ep.subscribers[address]; !exists {
		ep.subscribers[address] = []clients.Transaction{}
		return true
	}
	return false
}

func (ep *EthParser) GetTransactions(address string) []clients.Transaction {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	return ep.subscribers[strings.ToLower(address)]
}

func (ep *EthParser) FetchCurrentBlock() {
	payload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`

	rpcURL := os.Getenv("RPC_URL")
	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Failed to fetch current block: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	blockHex := result["result"].(string)
	var blockNumber int
	fmt.Sscanf(blockHex, "0x%x", &blockNumber)
	ep.mu.Lock()
	ep.currentBlock = blockNumber
	ep.mu.Unlock()
}

func (ep *EthParser) fetchTransactions(blockNumber int) []clients.Transaction {
	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x", true],"id":1}`, blockNumber)

	rpcURL := os.Getenv("RPC_URL")
	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Failed to fetch transactions: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	blockData, ok := result["result"].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid block data: %v", result["result"])
	}

	var transactions []clients.Transaction
	if transactionsData, ok := blockData["transactions"].([]interface{}); ok {
		for _, tx := range transactionsData {
			if txMap, ok := tx.(map[string]interface{}); ok {
				hash, _ := txMap["hash"].(string)
				from, _ := txMap["from"].(string)
				to, _ := txMap["to"].(string)
				value, _ := txMap["value"].(string)

				txs := clients.Transaction{
					Hash:  hash,
					From:  from,
					To:    to,
					Value: value,
					Block: blockNumber,
				}
				transactions = append(transactions, txs)
			} else {
				log.Printf("Unexpected transaction format: %v", tx)
			}
		}

	} else {
		log.Printf("No transactions found in block %d", blockNumber)
	}

	return transactions
}

func (ep *EthParser) ParseBlock(blockNumber int) {
	transactions := ep.fetchTransactions(blockNumber)
	ep.mu.Lock()
	for _, tx := range transactions {
		from := strings.ToLower(tx.From)
		to := strings.ToLower(tx.To)

		if _, exists := ep.subscribers[from]; exists {
			fmt.Println("User Exists in the from field")
			ep.subscribers[from] = append(ep.subscribers[from], tx)
		}
		if _, exists := ep.subscribers[to]; exists {
			fmt.Println("User Exists in the to field")
			ep.subscribers[to] = append(ep.subscribers[to], tx)
		}
	}
	ep.mu.Unlock()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
