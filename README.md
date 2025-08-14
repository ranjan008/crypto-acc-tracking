# Crypto Account Tracking

A comprehensive Ethereum wallet transaction tracker that fetches, categorizes, and exports transaction history to CSV format.

## Features

- 🔍 **Complete Transaction Tracking**: Fetches all transaction types from Ethereum wallets
  - Normal ETH transfers
  - Internal transactions
  - ERC-20 token transfers
  - ERC-721 NFT transfers
- 📊 **Automatic Categorization**: Intelligently categorizes transactions by type
- 💾 **CSV Export**: Exports structured data for portfolio management software
- ⚡ **High Performance**: Handles wallets with huge transactions
- 🛡️ **Robust**: Built-in retry logic and error handling
- 🎯 **Production Ready**: Rate limiting and validation

## Installation

### Prerequisites

- Go 1.21 or higher
- Internet connection for API access

### Build from Source

```bash
git clone <repository-url>
cd crypto-acc-tracking
go mod tidy
go build -o crypto-tracker
```

## Usage

### Basic Usage

```bash
# Track a wallet without API key (rate limited)
./crypto-tracker -a 0xa39b189482f984388a34460636fea9eb181ad1a6

# Track a wallet with Etherscan API key (recommended)
./crypto-tracker -a 0xa39b189482f984388a34460636fea9eb181ad1a6 -k YOUR_API_KEY

# Specify custom output file
./crypto-tracker -a 0xa39b189482f984388a34460636fea9eb181ad1a6 -o my_transactions.csv
```

### Command Line Options

- `-a, --address`: Ethereum wallet address to track (required)
- `-k, --api-key`: Etherscan API key (optional but recommended)
- `-o, --output`: Output CSV file path (default: transactions.csv)
- `-h, --help`: Show help information

### Getting an Etherscan API Key

1. Visit [https://etherscan.io/apis](https://etherscan.io/apis)
2. Create a free account
3. Generate an API key
4. Use the API key with the `-k` flag for higher rate limits

## Sample Addresses for Testing

- **Small wallet**: `0xa39b189482f984388a34460636fea9eb181ad1a6`
- **Medium wallet**: `0xd620AADaBaA20d2af700853C4504028cba7C3333`
- **Large wallet** (160k+ transactions): `0xfb50526f49894b78541b776f5aaefe43e3bd8590`

## CSV Output Format

The exported CSV file contains the following columns:

| Column | Description |
|--------|-------------|
| Transaction Hash | Unique identifier for the transaction |
| Date & Time | Transaction confirmation timestamp (UTC) |
| From Address | Sender's Ethereum address |
| To Address | Recipient's Ethereum address or contract |
| Transaction Type | ETH Transfer, ERC-20, ERC-721, Internal Transfer, Contract Interaction |
| Asset Contract Address | Contract address of the token or NFT (if applicable) |
| Asset Symbol / Name | Token symbol (e.g., ETH, USDC) or NFT collection name |
| Token ID | Unique identifier for NFTs (ERC-721, ERC-1155) |
| Value / Amount | Quantity of ETH or tokens transferred |
| Gas Fee (ETH) | Total transaction gas cost |
| Block Number | Block number where transaction was included |
| Status | Transaction status (Success, Failed, Unknown) |

## Architecture

```
crypto-acc-tracking/
├── cmd/                    # CLI command definitions
│   └── root.go
├── internal/
│   ├── etherscan/         # Etherscan API client
│   │   └── client.go
│   ├── models/            # Data structures
│   │   └── transaction.go
│   ├── processor/         # Transaction processing logic
│   │   └── processor.go
│   ├── exporter/          # CSV export functionality
│   │   └── csv.go
│   └── tracker/           # Main tracking logic
│       └── tracker.go
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

## Transaction Types

### Normal (External) Transfers
Direct ETH transfers between user-controlled addresses.

### Internal Transfers
Transfers that occur within smart contracts, not directly initiated by users.

### Token Transfers
- **ERC-20**: Fungible tokens (USDC, DAI, etc.)
- **ERC-721**: Non-fungible tokens (NFTs)

### Contract Interactions
Transactions that call smart contract functions.

## Performance Considerations

- **Rate Limiting**: Built-in delays between API calls to respect Etherscan limits
- **Pagination**: Handles large wallets by fetching data in chunks
- **Memory Efficient**: Processes transactions in batches
- **Retry Logic**: Automatic retry on API failures

## Error Handling

The application includes comprehensive error handling for:
- Invalid Ethereum addresses
- API rate limiting and failures
- Network connectivity issues
- Data parsing errors
- File system errors

## Disclaimer

This tool is for informational purposes only. Always verify transaction data independently before making financial decisions.
