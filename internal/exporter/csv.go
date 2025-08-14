package exporter

import (
	"crypto-acc-tracking/internal/models"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CSVExporter handles exporting transactions to CSV format
type CSVExporter struct {
	filename string
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter(filename string) *CSVExporter {
	return &CSVExporter{
		filename: filename,
	}
}

// Export writes transactions to a CSV file
func (e *CSVExporter) Export(transactions []*models.Transaction) error {
	file, err := os.Create(e.filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"Transaction Hash",
		"Date & Time",
		"From Address",
		"To Address",
		"Transaction Type",
		"Asset Contract Address",
		"Asset Symbol / Name",
		"Token ID",
		"Value / Amount",
		"Gas Fee (ETH)",
		"Block Number",
		"Status",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write transaction data
	for _, tx := range transactions {
		record := []string{
			tx.Hash,
			tx.DateTime.Format("2006-01-02 15:04:05 UTC"),
			tx.FromAddress,
			tx.ToAddress,
			string(tx.TransactionType),
			tx.AssetContractAddr,
			e.formatAssetInfo(tx.AssetSymbol, tx.AssetName),
			tx.TokenID,
			e.formatValue(tx.ValueFormatted, tx.AssetSymbol),
			tx.GasFeeETH,
			tx.BlockNumber,
			tx.Status,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// formatAssetInfo combines symbol and name for better readability
func (e *CSVExporter) formatAssetInfo(symbol, name string) string {
	if symbol == "" && name == "" {
		return ""
	}
	if symbol == "" {
		return name
	}
	if name == "" {
		return symbol
	}
	if symbol == name {
		return symbol
	}
	return fmt.Sprintf("%s (%s)", symbol, name)
}

// formatValue formats the value with symbol for better readability
func (e *CSVExporter) formatValue(value, symbol string) string {
	if value == "" || value == "0" {
		return "0"
	}

	// Remove trailing zeros for cleaner display
	value = strings.TrimRight(value, "0")
	value = strings.TrimRight(value, ".")

	if symbol == "" {
		return value
	}

	return fmt.Sprintf("%s %s", value, symbol)
}

// GetExportSummary returns a summary of the export operation
func (e *CSVExporter) GetExportSummary(transactions []*models.Transaction) map[string]interface{} {
	summary := map[string]interface{}{
		"total_transactions": len(transactions),
		"filename":           e.filename,
	}

	// Count transactions by type
	typeCounts := make(map[models.TransactionType]int)
	for _, tx := range transactions {
		typeCounts[tx.TransactionType]++
	}

	summary["transaction_types"] = typeCounts

	// Count unique assets
	assets := make(map[string]bool)
	for _, tx := range transactions {
		if tx.AssetSymbol != "" {
			assets[tx.AssetSymbol] = true
		}
	}
	summary["unique_assets"] = len(assets)

	return summary
}
