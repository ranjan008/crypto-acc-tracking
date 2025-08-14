package processor

import (
	"crypto-acc-tracking/internal/models"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	WeiPerEth = 1e18
)

// Processor handles transaction data processing and conversion
type Processor struct{}

// New creates a new processor instance
func New() *Processor {
	return &Processor{}
}

// ProcessNormalTransaction converts Etherscan normal transaction to unified format
func (p *Processor) ProcessNormalTransaction(tx models.EtherscanNormalTx) (*models.Transaction, error) {
	timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	value, ok := new(big.Int).SetString(tx.Value, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse value: %s", tx.Value)
	}

	gasFee := p.calculateGasFee(tx.GasUsed, tx.GasPrice)
	valueFormatted := p.formatEthValue(value)

	transaction := &models.Transaction{
		Hash:              tx.Hash,
		DateTime:          time.Unix(timestamp, 0).UTC(),
		FromAddress:       tx.From,
		ToAddress:         tx.To,
		TransactionType:   models.ETHTransfer,
		AssetContractAddr: "", // ETH doesn't have a contract address
		AssetSymbol:       "ETH",
		AssetName:         "Ethereum",
		TokenID:           "",
		Value:             value,
		ValueFormatted:    valueFormatted,
		GasFeeETH:         gasFee,
		BlockNumber:       tx.BlockNumber,
		TransactionIndex:  tx.TransactionIndex,
		Status:            p.getTransactionStatus(tx.IsError, tx.TxReceiptStatus),
	}

	// Check if it's a contract interaction
	if tx.To != "" && tx.Input != "0x" && tx.Input != "" {
		transaction.TransactionType = models.ContractCall
	}

	return transaction, nil
}

// ProcessInternalTransaction converts Etherscan internal transaction to unified format
func (p *Processor) ProcessInternalTransaction(tx models.EtherscanInternalTx) (*models.Transaction, error) {
	timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	value, ok := new(big.Int).SetString(tx.Value, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse value: %s", tx.Value)
	}

	valueFormatted := p.formatEthValue(value)

	transaction := &models.Transaction{
		Hash:              tx.Hash,
		DateTime:          time.Unix(timestamp, 0).UTC(),
		FromAddress:       tx.From,
		ToAddress:         tx.To,
		TransactionType:   models.InternalTx,
		AssetContractAddr: tx.ContractAddress,
		AssetSymbol:       "ETH",
		AssetName:         "Ethereum",
		TokenID:           "",
		Value:             value,
		ValueFormatted:    valueFormatted,
		GasFeeETH:         "0", // Internal transactions don't have gas fees
		BlockNumber:       tx.BlockNumber,
		TransactionIndex:  "",
		Status:            p.getInternalTransactionStatus(tx.IsError),
	}

	return transaction, nil
}

// ProcessTokenTransaction converts Etherscan token transaction to unified format
func (p *Processor) ProcessTokenTransaction(tx models.EtherscanTokenTx) (*models.Transaction, error) {
	timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	value, ok := new(big.Int).SetString(tx.Value, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse value: %s", tx.Value)
	}

	gasFee := p.calculateGasFee(tx.GasUsed, tx.GasPrice)

	// Parse token decimals
	decimals, err := strconv.Atoi(tx.TokenDecimal)
	if err != nil {
		decimals = 18 // Default to 18 if parsing fails
	}

	valueFormatted := p.formatTokenValue(value, decimals)

	transaction := &models.Transaction{
		Hash:              tx.Hash,
		DateTime:          time.Unix(timestamp, 0).UTC(),
		FromAddress:       tx.From,
		ToAddress:         tx.To,
		TransactionType:   models.ERC20Transfer,
		AssetContractAddr: tx.ContractAddress,
		AssetSymbol:       tx.TokenSymbol,
		AssetName:         tx.TokenName,
		TokenID:           "",
		Value:             value,
		ValueFormatted:    valueFormatted,
		GasFeeETH:         gasFee,
		BlockNumber:       tx.BlockNumber,
		TransactionIndex:  tx.TransactionIndex,
		Status:            "1", // Token transactions are usually successful if they appear in the list
	}

	return transaction, nil
}

// ProcessNFTTransaction converts Etherscan NFT transaction to unified format
func (p *Processor) ProcessNFTTransaction(tx models.EtherscanNFTTx) (*models.Transaction, error) {
	timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	gasFee := p.calculateGasFee(tx.GasUsed, tx.GasPrice)

	transaction := &models.Transaction{
		Hash:              tx.Hash,
		DateTime:          time.Unix(timestamp, 0).UTC(),
		FromAddress:       tx.From,
		ToAddress:         tx.To,
		TransactionType:   models.ERC721Transfer,
		AssetContractAddr: tx.ContractAddress,
		AssetSymbol:       tx.TokenSymbol,
		AssetName:         tx.TokenName,
		TokenID:           tx.TokenID,
		Value:             big.NewInt(1), // NFTs typically have quantity of 1
		ValueFormatted:    "1",
		GasFeeETH:         gasFee,
		BlockNumber:       tx.BlockNumber,
		TransactionIndex:  tx.TransactionIndex,
		Status:            "1", // NFT transactions are usually successful if they appear in the list
	}

	return transaction, nil
}

// calculateGasFee calculates the gas fee in ETH
func (p *Processor) calculateGasFee(gasUsed, gasPrice string) string {
	gasUsedBig, ok1 := new(big.Int).SetString(gasUsed, 10)
	gasPriceBig, ok2 := new(big.Int).SetString(gasPrice, 10)

	if !ok1 || !ok2 {
		return "0"
	}

	gasFeeWei := new(big.Int).Mul(gasUsedBig, gasPriceBig)
	return p.formatEthValue(gasFeeWei)
}

// formatEthValue formats wei value to ETH with proper decimal places
func (p *Processor) formatEthValue(weiValue *big.Int) string {
	if weiValue.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	ethValue := new(big.Float).SetInt(weiValue)
	divisor := new(big.Float).SetFloat64(WeiPerEth)
	ethValue.Quo(ethValue, divisor)

	return ethValue.Text('f', 18)
}

// formatTokenValue formats token value with proper decimal places
func (p *Processor) formatTokenValue(value *big.Int, decimals int) string {
	if value.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	if decimals == 0 {
		return value.String()
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	tokenValue := new(big.Float).SetInt(value)
	divisorFloat := new(big.Float).SetInt(divisor)
	tokenValue.Quo(tokenValue, divisorFloat)

	return tokenValue.Text('f', int(decimals))
}

// getTransactionStatus determines transaction status
func (p *Processor) getTransactionStatus(isError, txReceiptStatus string) string {
	if isError == "1" {
		return "Failed"
	}
	if txReceiptStatus == "1" {
		return "Success"
	}
	if txReceiptStatus == "0" {
		return "Failed"
	}
	return "Unknown"
}

// getInternalTransactionStatus determines internal transaction status
func (p *Processor) getInternalTransactionStatus(isError string) string {
	if isError == "1" {
		return "Failed"
	}
	return "Success"
}

// SortTransactionsByTime sorts transactions by timestamp in descending order
func (p *Processor) SortTransactionsByTime(transactions []*models.Transaction) {
	for i := 0; i < len(transactions)-1; i++ {
		for j := 0; j < len(transactions)-i-1; j++ {
			if transactions[j].DateTime.Before(transactions[j+1].DateTime) {
				transactions[j], transactions[j+1] = transactions[j+1], transactions[j]
			}
		}
	}
}

// DeduplicateTransactions removes duplicate transactions based on hash and type
func (p *Processor) DeduplicateTransactions(transactions []*models.Transaction) []*models.Transaction {
	seen := make(map[string]bool)
	var unique []*models.Transaction

	for _, tx := range transactions {
		key := fmt.Sprintf("%s_%s", tx.Hash, tx.TransactionType)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, tx)
		}
	}

	return unique
}

// ValidateEthereumAddress validates if the provided string is a valid Ethereum address
func (p *Processor) ValidateEthereumAddress(address string) bool {
	if len(address) != 42 {
		return false
	}

	if !strings.HasPrefix(address, "0x") && !strings.HasPrefix(address, "0X") {
		return false
	}

	// Check if all characters after 0x are hexadecimal
	hex := address[2:]
	for _, char := range hex {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}
