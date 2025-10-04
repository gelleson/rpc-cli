# Default config (no label)
config {
  url = "https://eth-mainnet.g.alchemy.com/v2/demo"
  headers = {
    Content-Type = "application/json"
  }
  timeout = 30
}

# Production config
config "production" {
  url = "https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY"
  headers = {
    Content-Type  = "application/json"
    Authorization = "Bearer prod_token_secret_key_12345"
  }
  timeout = 60
}

# Staging config
config "staging" {
  url = "https://eth-sepolia.g.alchemy.com/v2/demo"
  headers = {
    Content-Type = "application/json"
    X-Environment = "staging"
  }
  timeout = 45
}

# Simple request - get latest block number
request "get_block_number" {
  method = "eth_blockNumber"
  params = []
}

# Request with simple parameters
request "get_balance" {
  method = "eth_getBalance"
  params = ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]
}

# Request with config override
request "get_balance_staging" {
  config = "staging"
  method = "eth_getBalance"
  params = ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]
}

# Request with URL override
request "get_balance_custom_rpc" {
  url = "https://cloudflare-eth.com"
  method = "eth_getBalance"
  params = ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "latest"]
}

# Request with header overrides
request "admin_call" {
  headers = {
    Authorization = "Bearer admin_token_xyz"
    X-Admin-Key   = "secret123"
  }
  method = "eth_blockNumber"
  params = []
}

# Request with timeout override
request "long_query" {
  timeout = 120
  method = "eth_getLogs"
  params = [
    {
      fromBlock = "0x1000000"
      toBlock   = "0x1000100"
    }
  ]
}

# Complex nested params - get transaction by hash
request "get_transaction" {
  method = "eth_getTransactionByHash"
  params = ["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"]
}

# Request with object parameters
request "get_logs" {
  method = "eth_getLogs"
  params = [
    {
      fromBlock = "latest"
      toBlock   = "latest"
      address   = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
      topics    = []
    }
  ]
}

# Complex nested structures with multiple objects
request "batch_transfer" {
  method = "wallet.batchTransfer"
  params = [
    {
      to     = "0x123456789abcdef123456789abcdef123456789a"
      amount = 100
      token  = "USDT"
    },
    {
      to     = "0x456789abcdef123456789abcdef123456789abcde"
      amount = 250
      token  = "USDC"
    },
    {
      to     = "0x789abcdef123456789abcdef123456789abcdef12"
      amount = 500
      token  = "DAI"
    }
  ]
}

# Deeply nested structures
request "complex_query" {
  method = "db.query"
  params = [
    "users",
    {
      age = {
        gt = 18
        lt = 65
      }
      status = "active"
      tags   = ["premium", "verified"]
    },
    {
      limit  = 100
      offset = 0
      sort = {
        field = "created_at"
        order = "desc"
      }
    }
  ]
}

# Request using production config
request "get_balance_prod" {
  config = "production"
  method = "eth_getBalance"
  params = ["0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", "latest"]
}

# Request with boolean and numeric values
request "complex_params" {
  method = "eth_call"
  params = [
    {
      to   = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
      data = "0x70a08231000000000000000000000000742d35cc6634c0532925a3b844bc9e7595f0beb"
    },
    "latest"
  ]
}

# Request demonstrating array parameters
request "get_multiple_blocks" {
  method = "eth_getBlockByNumber"
  params = ["0x1000000", true]
}
