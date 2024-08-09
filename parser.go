package main

import (
	"context"
	"sync"
	"log"
	"encoding/json"
	"fmt"
	"math/big"
)

// Chainstack Endpoint
const NodeURL = "https://nd-000-364-211.p2pify.com/5b8d22690a57f293b3a1ed8758014e35"


// Parser Interface
type Parser interface {
	// last parsed block
	GetCurrentBlock() *big.Int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

// EthParser implementation of Parser for Ethereum blockchain.
type EthParser struct {
	subscribedAddresses map[string]bool
	mu                  sync.Mutex
	notificationService *NotificationService
	ctx                 context.Context
	cancel              context.CancelFunc
	transactions map[string][]Transaction
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	TransactionIndex string `json:"transactionIndex"`
}

// Block structure to hold block data
type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

// NewEthParser creates a new EthParser.
func NewEthParser(notificationService *NotificationService) *EthParser {
	ctx, cancel := context.WithCancel(context.Background())
	return &EthParser{
		subscribedAddresses: make(map[string]bool),
		notificationService: notificationService,
		ctx:                 ctx,
		cancel:              cancel,
	}
}

// // SubscribeAddress subscribes an address for transaction monitoring.
func (p *EthParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	hex := AddressToHex(address)
	p.subscribedAddresses[hex] = true

	return p.subscribedAddresses[hex]
}

func (p *EthParser) GetCurrentBlock() *big.Int {
	response, err := SendRPCRequest(NodeURL, "eth_blockNumber", []interface{}{})
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	var blockNumberHex string
	if err := json.Unmarshal(response.Result, &blockNumberHex); err != nil {
		log.Fatalf("Failed to parse block number: %v", err)
	}

	blockNumber, err := HexToBigInt(blockNumberHex)
	if err != nil {
		log.Fatalf("Failed to convert block number: %v", err)
	}

	return blockNumber
}

func (p *EthParser) GetTransactions(address string) []Transaction {
	if transactions, found := p.transactions[address]; found {
		return transactions
	} else {
		return make([]Transaction, 0)
	}
}

// GetBlockByNumber retrieves a block by its number
func (p *EthParser) GetBlockByNumber(blockNumber string) (*Block, error) {
	// Fetch the block by number
	response, err := SendRPCRequest(NodeURL, "eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %v", err)
	}

	// Parse the block response
	var block Block
	if err := json.Unmarshal(response.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to parse block: %v", err)
	}

	return &block, nil
}

// processBlock processes the transactions in a block.
func (p *EthParser) processBlock(block *Block) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, tx := range block.Transactions {
		from := AddressToHex(tx.From)
		to := AddressToHex(tx.To)
		p.AppendTransaction(from, tx)
		p.AppendTransaction(to, tx)

		// Check if the sender or receiver address is subscribed.
		if p.isSubscribed(from) {
			p.notificationService.Notify(from, tx, false)
		}

		if p.isSubscribed(to) {
			p.notificationService.Notify(to, tx, true)
		}
	}
}

// isSubscribed checks if an address is subscribed.
func (p *EthParser) isSubscribed(address string) bool {
	return p.subscribedAddresses[address]
}

func (p *EthParser) GetTransactionsForBlock(blockNumber string) {
	block, err := p.GetBlockByNumber(blockNumber)
	if err != nil {
		fmt.Printf("Error In getting BlockByNumber %v", err)
	}

	p.processBlock(block)
}

func (p *EthParser) AppendTransaction(address string, txn Transaction) {
	if transactions, found := p.transactions[address]; found {
		p.transactions[address] = append(transactions, txn)
	} else {
		p.transactions[address] = make([]Transaction, 0)
	}
}
