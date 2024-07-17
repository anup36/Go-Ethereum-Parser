# Ethereum Blockchain Parser Assignment Using Go Lang

## Usage :
This project uses Go Lang to expose Parser Apis which can be called to:
    - subscribe to all transactions on chain for an ethereum address.
    - get all transactions for the subscribed address - both incoming and outgoing transactions.
    - get current block from the ethereum blockchain.

### API Endpoints:

**1. Get Current Block**
```bash
curl http://localhost:8080/current-block
```

**2. Get Transactions for Address**
```bash
curl http://localhost:8080/transactions\?address\={your_ethereum_address}
```

**3. Subscribe to an Address**
```bash
curl http://localhost:8080/subscribe\?address\={your_ethereum_address}
```

## Run Application:
```bash
go run main.go
```