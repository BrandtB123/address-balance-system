package db

import (
	"fmt"
	"strconv"
	"unit410/models"

	"github.com/go-pg/pg"
)

var DB *pg.DB

func NewDB() {
	db := pg.Connect(&pg.Options{
		Addr:     "unit410.cs96vtk3gnbe.us-east-2.rds.amazonaws.com:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "unit410",
	})
	DB = db
}

func CreateAddressesTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Addresses (
			UUID              UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			Address           TEXT,
			Network           TEXT,
			SignificantDigits INT,
			Asset             TEXT,
			CONSTRAINT uc_address_asset UNIQUE (Address, Asset)
		)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create Addresses table: %v", err)
	}

	return nil
}

func CreateBalancesTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Balances (
			UUID    UUID REFERENCES Addresses(UUID),
			Timestamp          TIMESTAMP,
			Balance            TEXT,
			Available          TEXT,
			DelegatableVesting TEXT,
			Delegated          TEXT,
			Staking            TEXT,
			Unbonding          TEXT,
			Reward             TEXT
		)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create Balances table: %v", err)
	}

	return nil
}

func AddAddress(address models.Address) error {
	query := `
		INSERT INTO Addresses (Address, Network, SignificantDigits, Asset)
		VALUES ('%s', '%s', '%s', '%s')
		ON CONFLICT (Address, Asset) DO NOTHING
	`
	q := fmt.Sprintf(query, address.Address, address.Network, strconv.Itoa(address.SignificantDigits), address.Asset)

	fmt.Println(q)
	_, err := DB.Exec(q)
	if err != nil {
		return fmt.Errorf("failed to add address: %v", err)
	}

	return nil
}

func AddBalance(asset string, balance models.Bal) error {
	query := `
		INSERT INTO Balances (UUID, Timestamp, Balance, Available, DelegatableVesting, Delegated, Staking, Unbonding, Reward)
		VALUES ((SELECT UUID FROM Addresses WHERE Address = '%s' AND Asset = '%s'), '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')
	`
	s := balance.Timestamp.Format("2006-01-02 15:04:05")

	q := fmt.Sprintf(query, balance.Address, asset, s, balance.Balance, balance.Available, balance.DelegatableVesting, balance.Delegated, balance.Staking, balance.Unbonding, balance.Reward)

	_, err := DB.Exec(q)
	if err != nil {
		return fmt.Errorf("failed to add balance: %v", err)
	}

	return nil
}

// func SaveBlockData(balance models.Bal) error {
// 	_, err := DB.Model(&block).OnConflict("(height) DO UPDATE").Set("tx_count = ?tx_count").Insert()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func GetLatestHeight() (int64, error) {
// 	var block Block
// 	err := DB.Model(&block).Order("height DESC").Limit(1).Select()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to retrieve the latest block height: %v", err)
// 	}

// 	return block.Height, nil
// }

// func GetTransactionsInPastNBlocks(n int) (int, error) {
// 	var block Block
// 	var totalTxs int
// 	err := DB.Model(&block).ColumnExpr("SUM(tx_count)").Limit(n).Select(&totalTxs)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return totalTxs, nil
// }

// func GetProposedBlocksByValidator(proposer string) ([]int64, error) {
// 	var heights []int64
// 	_, err := DB.Query(&heights, `
// 		SELECT height
// 		FROM blocks
// 		WHERE proposer = ?
// 	`, proposer)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return heights, nil
// }

// func GetTopNPeersByScore(n int) ([]PeerScore, error) {
// 	var blocks []Block
// 	err := DB.Model(&blocks).
// 		Order("height DESC").
// 		Limit(n).
// 		Select()

// 	if err != nil {
// 		return nil, err
// 	}

// 	peers := make(map[string]int)
// 	for _, block := range blocks {
// 		for _, score := range block.Peers {
// 			peers[score.Address] += score.Score
// 		}
// 	}

// 	var topPeers []PeerScore

// 	for address, score := range peers {
// 		topPeers = append(topPeers, PeerScore{
// 			Address: address,
// 			Score:   score,
// 		})
// 	}

// 	sort.Slice(topPeers, func(i, j int) bool {
// 		return topPeers[i].Score > topPeers[j].Score
// 	})

// 	if len(topPeers) > n {
// 		topPeers = topPeers[:n]
// 	}
// 	return topPeers, nil
// }
