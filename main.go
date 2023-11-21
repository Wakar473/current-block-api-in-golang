// package main

// import (
// 	"context"
// 	"encoding/json"
// 	// "fmt"
// 	"github.com/ethereum/go-ethereum/rpc"
// 	"github.com/gin-gonic/gin"
// 	"log"
// 	// "math/big"
// )

// type BlockDetails struct {
// 	Number        string `json:"number"`
// 	ParentHash    string `json:"parentHash"`
// 	BlockHash     string `json:"hash"`
// 	Timestamp     string `json:"timestamp"`
// 	Transactions  int    `json:"transactions"`
// }

// func main() {
// 	r := gin.Default()

// 	ethEndpoint := "https://rpc-alpha-testnet.saitascan.io"

// 	client, err := rpc.Dial(ethEndpoint)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
// 	}

// 	defer client.Close()

// 	r.GET("/block/:blockNumber", func(c *gin.Context) {
// 		blockNumber := c.Param("blockNumber")

// 		var block map[string]interface{}
// 		if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", "0x"+blockNumber, true); err != nil {
// 			log.Printf("Failed to retrieve block details for block number %s: %v", blockNumber, err)
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		blockDetails := &BlockDetails{
// 			Number:        block["number"].(string),
// 			ParentHash:    block["parentHash"].(string),
// 			BlockHash:     block["hash"].(string),
// 			Timestamp:     block["timestamp"].(string),
// 			Transactions: len(block["transactions"].([]interface{})),
// 		}

// 		jsonData, err := json.Marshal(blockDetails)
// 		if err != nil {
// 			log.Printf("Failed to marshal block details to JSON: %v", err)
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(200, string(jsonData))
// 	})

// 	r.Run(":8080")
// }

package main

import (
	"context"
	// "encoding/json"
	// "fmt"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"

	// "math/big"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type BlockDetails struct {
	Number       string `json:"number"`
	ParentHash   string `json:"parentHash"`
	BlockHash    string `json:"hash"`
	Timestamp    string `json:"timestamp"`
	Transactions int    `json:"transactions"`
}

func main() {
	ethEndpoint := "https://rpc-alpha-testnet.saitascan.io"

	client, err := rpc.Dial(ethEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	defer client.Close()

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/saita")
	if err != nil {
		log.Fatalf("Failed to connect to the MySQL database: %v", err)
	}

	defer db.Close()

	r := gin.Default()

	r.GET("/block", func(c *gin.Context) {
		// Fetch the latest block number
		var latestBlockNum string
		if err := client.CallContext(context.Background(), &latestBlockNum, "eth_blockNumber"); err != nil {
			log.Printf("Failed to retrieve latest block number: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Retrieve block details for the current block
		var block map[string]interface{}
		if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", latestBlockNum, true); err != nil {
			log.Printf("Failed to retrieve block details for block number %s: %v", latestBlockNum, err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		blockDetails := &BlockDetails{
			Number:       block["number"].(string),
			ParentHash:   block["parentHash"].(string),
			BlockHash:    block["hash"].(string),
			Timestamp:    block["timestamp"].(string),
			Transactions: len(block["transactions"].([]interface{})),
		}

		// Store block details in MySQL database
		stmt, err := db.Prepare(`INSERT INTO block_details (block_number, parent_hash, block_hash, timestamp, transaction_count) VALUES (?, ?, ?, ?, ?)`)
		if err != nil {
			log.Printf("Failed to prepare SQL statement: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(blockDetails.Number, blockDetails.ParentHash, blockDetails.BlockHash, blockDetails.Timestamp, blockDetails.Transactions)
		if err != nil {
			log.Printf("Failed to insert block details into MySQL: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, blockDetails)
	})

	// Start the Gin router
	r.Run(":8080")

	// Schedule periodic updates to fetch and store block details
	ticker := time.NewTicker(time.Minute * 10) // Update every 10 minutes
	defer ticker.Stop()

	for {
		<-ticker.C
		log.Println("Fetching and storing current block details...")

		// Fetch and store block details
		// ... (Same as above)
	}
}
