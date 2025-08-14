package models

import (
	"math/big"
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	ETHTransfer     TransactionType = "ETH Transfer"
	ERC20Transfer   TransactionType = "ERC-20 Transfer"
	ERC721Transfer  TransactionType = "ERC-721 Transfer"
	ERC1155Transfer TransactionType = "ERC-1155 Transfer"
	ContractCall    TransactionType = "Contract Interaction"
	InternalTx      TransactionType = "Internal Transfer"
)

// Transaction represents a unified transaction structure
type Transaction struct {
	Hash              string          `json:"hash"`
	DateTime          time.Time       `json:"dateTime"`
	FromAddress       string          `json:"fromAddress"`
	ToAddress         string          `json:"toAddress"`
	TransactionType   TransactionType `json:"transactionType"`
	AssetContractAddr string          `json:"assetContractAddr"`
	AssetSymbol       string          `json:"assetSymbol"`
	AssetName         string          `json:"assetName"`
	TokenID           string          `json:"tokenId"`
	Value             *big.Int        `json:"value"`
	ValueFormatted    string          `json:"valueFormatted"`
	GasFeeETH         string          `json:"gasFeeEth"`
	BlockNumber       string          `json:"blockNumber"`
	TransactionIndex  string          `json:"transactionIndex"`
	Status            string          `json:"status"`
}

// EtherscanNormalTx represents a normal transaction from Etherscan API
type EtherscanNormalTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
	MethodID          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
}

// EtherscanInternalTx represents an internal transaction from Etherscan API
type EtherscanInternalTx struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Input           string `json:"input"`
	Type            string `json:"type"`
	Gas             string `json:"gas"`
	GasUsed         string `json:"gasUsed"`
	TraceID         string `json:"traceId"`
	IsError         string `json:"isError"`
	ErrCode         string `json:"errCode"`
}

// EtherscanTokenTx represents a token transfer from Etherscan API
type EtherscanTokenTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`
}

// EtherscanNFTTx represents an NFT transfer from Etherscan API
type EtherscanNFTTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	From              string `json:"from"`
	To                string `json:"to"`
	TokenID           string `json:"tokenID"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	Confirmations     string `json:"confirmations"`
}

// EtherscanResponse represents the API response structure
type EtherscanResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
