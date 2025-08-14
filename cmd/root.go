package cmd

import (
	"crypto-acc-tracking/internal/tracker"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	address string
	apiKey  string
	output  string
)

var rootCmd = &cobra.Command{
	Use:   "crypto-tracker",
	Short: "Ethereum wallet transaction tracker",
	Long:  `A CLI tool to track and export Ethereum wallet transactions to CSV format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if address == "" {
			return fmt.Errorf("ethereum address is required")
		}

		t := tracker.New(apiKey)
		return t.TrackWallet(address, output)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&address, "address", "a", "", "Ethereum wallet address to track (required)")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Etherscan API key (optional but recommended for higher rate limits)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "transactions.csv", "Output CSV file path")

	rootCmd.MarkFlagRequired("address")
}
