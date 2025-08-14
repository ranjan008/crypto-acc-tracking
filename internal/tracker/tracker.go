package tracker

import (
	"crypto-acc-tracking/internal/etherscan"
	"crypto-acc-tracking/internal/exporter"
	"crypto-acc-tracking/internal/models"
	"crypto-acc-tracking/internal/processor"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultPageSize = 10000 // Max transactions per request
	MaxPages        = 100   // Limit pages to prevent infinite loops
)

// Tracker represents the main transaction tracking service
type Tracker struct {
	etherscanClient *etherscan.Client
	processor       *processor.Processor
}

// New creates a new tracker instance
func New(apiKey string) *Tracker {
	return &Tracker{
		etherscanClient: etherscan.New(apiKey),
		processor:       processor.New(),
	}
}

// TrackWallet retrieves and exports all transactions for a wallet address
func (t *Tracker) TrackWallet(address, outputFile string) error {
	// Validate address
	if !t.processor.ValidateEthereumAddress(address) {
		return fmt.Errorf("invalid Ethereum address: %s", address)
	}

	// Normalize address to lowercase
	address = strings.ToLower(address)

	fmt.Printf("🔍 Tracking wallet: %s\n", address)
	fmt.Printf("📊 Fetching transaction data...\n\n")

	// Fetch all transaction types
	var allTransactions []*models.Transaction

	// 1. Fetch normal transactions
	fmt.Printf("⏳ Fetching normal transactions...\n")
	normalTxs, err := t.fetchAllNormalTransactions(address)
	if err != nil {
		return fmt.Errorf("failed to fetch normal transactions: %w", err)
	}
	fmt.Printf("✅ Found %d normal transactions\n", len(normalTxs))

	// 2. Fetch internal transactions
	fmt.Printf("⏳ Fetching internal transactions...\n")
	internalTxs, err := t.fetchAllInternalTransactions(address)
	if err != nil {
		return fmt.Errorf("failed to fetch internal transactions: %w", err)
	}
	fmt.Printf("✅ Found %d internal transactions\n", len(internalTxs))

	// Wait before next API call batch
	time.Sleep(2 * time.Second)

	// 3. Fetch token transactions
	fmt.Printf("⏳ Fetching ERC-20 token transactions...\n")
	tokenTxs, err := t.fetchAllTokenTransactions(address)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to fetch token transactions: %v\n", err)
		fmt.Printf("📄 Continuing with available data...\n")
		tokenTxs = []*models.Transaction{} // Use empty slice
	} else {
		fmt.Printf("✅ Found %d token transactions\n", len(tokenTxs))
	}

	// Wait before next API call batch
	time.Sleep(2 * time.Second)

	// 4. Fetch NFT transactions
	fmt.Printf("⏳ Fetching ERC-721 NFT transactions...\n")
	nftTxs, err := t.fetchAllNFTTransactions(address)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to fetch NFT transactions: %v\n", err)
		fmt.Printf("📄 Continuing with available data...\n")
		nftTxs = []*models.Transaction{} // Use empty slice
	} else {
		fmt.Printf("✅ Found %d NFT transactions\n", len(nftTxs))
	}

	// Combine all transactions
	allTransactions = append(allTransactions, normalTxs...)
	allTransactions = append(allTransactions, internalTxs...)
	allTransactions = append(allTransactions, tokenTxs...)
	allTransactions = append(allTransactions, nftTxs...)

	// Deduplicate and sort
	fmt.Printf("\n📋 Processing transactions...\n")
	allTransactions = t.processor.DeduplicateTransactions(allTransactions)
	t.processor.SortTransactionsByTime(allTransactions)

	fmt.Printf("✅ Total unique transactions: %d\n", len(allTransactions))

	// Export to CSV
	fmt.Printf("\n💾 Exporting to CSV: %s\n", outputFile)
	csvExporter := exporter.NewCSVExporter(outputFile)
	if err := csvExporter.Export(allTransactions); err != nil {
		return fmt.Errorf("failed to export to CSV: %w", err)
	}

	// Print summary
	summary := csvExporter.GetExportSummary(allTransactions)
	t.printSummary(summary)

	fmt.Printf("\n🎉 Export completed successfully!\n")
	return nil
}

// fetchAllNormalTransactions fetches all normal transactions with pagination
func (t *Tracker) fetchAllNormalTransactions(address string) ([]*models.Transaction, error) {
	var allTransactions []*models.Transaction
	page := 1

	for page <= MaxPages {
		txs, err := t.etherscanClient.GetNormalTransactions(address, 0, 99999999, page, DefaultPageSize)
		if err != nil {
			return nil, err
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			processedTx, err := t.processor.ProcessNormalTransaction(tx)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to process normal transaction %s: %v\n", tx.Hash, err)
				continue
			}
			allTransactions = append(allTransactions, processedTx)
		}

		if len(txs) < DefaultPageSize {
			break
		}

		page++
		time.Sleep(3 * time.Second) // Rate limiting
	}

	return allTransactions, nil
}

// fetchAllInternalTransactions fetches all internal transactions with pagination
func (t *Tracker) fetchAllInternalTransactions(address string) ([]*models.Transaction, error) {
	var allTransactions []*models.Transaction
	page := 1

	for page <= MaxPages {
		txs, err := t.etherscanClient.GetInternalTransactions(address, 0, 99999999, page, DefaultPageSize)
		if err != nil {
			return nil, err
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			processedTx, err := t.processor.ProcessInternalTransaction(tx)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to process internal transaction %s: %v\n", tx.Hash, err)
				continue
			}
			allTransactions = append(allTransactions, processedTx)
		}

		if len(txs) < DefaultPageSize {
			break
		}

		page++
		time.Sleep(3 * time.Second) // Rate limiting
	}

	return allTransactions, nil
}

// fetchAllTokenTransactions fetches all token transactions with pagination
func (t *Tracker) fetchAllTokenTransactions(address string) ([]*models.Transaction, error) {
	var allTransactions []*models.Transaction
	page := 1

	for page <= MaxPages {
		txs, err := t.etherscanClient.GetTokenTransactions(address, 0, 99999999, page, DefaultPageSize)
		if err != nil {
			return nil, err
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			processedTx, err := t.processor.ProcessTokenTransaction(tx)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to process token transaction %s: %v\n", tx.Hash, err)
				continue
			}
			allTransactions = append(allTransactions, processedTx)
		}

		if len(txs) < DefaultPageSize {
			break
		}

		page++
		time.Sleep(3 * time.Second) // Rate limiting
	}

	return allTransactions, nil
}

// fetchAllNFTTransactions fetches all NFT transactions with pagination
func (t *Tracker) fetchAllNFTTransactions(address string) ([]*models.Transaction, error) {
	var allTransactions []*models.Transaction
	page := 1

	for page <= MaxPages {
		txs, err := t.etherscanClient.GetNFTTransactions(address, 0, 99999999, page, DefaultPageSize)
		if err != nil {
			return nil, err
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			processedTx, err := t.processor.ProcessNFTTransaction(tx)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to process NFT transaction %s: %v\n", tx.Hash, err)
				continue
			}
			allTransactions = append(allTransactions, processedTx)
		}

		if len(txs) < DefaultPageSize {
			break
		}

		page++
		time.Sleep(3 * time.Second) // Rate limiting
	}

	return allTransactions, nil
}

// printSummary prints a summary of the export operation
func (t *Tracker) printSummary(summary map[string]interface{}) {
	fmt.Printf("\n📈 Export Summary:\n")
	fmt.Printf("   Total Transactions: %d\n", summary["total_transactions"])
	fmt.Printf("   Unique Assets: %d\n", summary["unique_assets"])
	fmt.Printf("   Output File: %s\n", summary["filename"])

	if typeCounts, ok := summary["transaction_types"].(map[models.TransactionType]int); ok {
		fmt.Printf("\n📊 Transaction Types:\n")
		for txType, count := range typeCounts {
			fmt.Printf("   %s: %d\n", txType, count)
		}
	}
}
