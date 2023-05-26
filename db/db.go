package db

import (
	"fmt"
	"strconv"
	"time"
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
			Timestamp          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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
		INSERT INTO Balances (UUID, Balance, Available, DelegatableVesting, Delegated, Staking, Unbonding, Reward)
		VALUES ((SELECT UUID FROM Addresses WHERE Address = '%s' AND Asset = '%s'), '%s', '%s', '%s', '%s', '%s', '%s', '%s')
	`
	q := fmt.Sprintf(query, balance.Address, asset, balance.Balance, balance.Available, balance.DelegatableVesting, balance.Delegated, balance.Staking, balance.Unbonding, balance.Reward)

	_, err := DB.Exec(q)
	if err != nil {
		return fmt.Errorf("failed to add balance: %v", err)
	}

	return nil
}

type BalanceAsset struct {
	UUID               string    `pg:"uuid"`
	Timestamp          time.Time `pg:"timestamp"`
	Balance            string    `pg:"balance"`
	Available          string    `pg:"available"`
	Delegatablevesting string    `pg:"delegatablevesting"`
	Delegated          string    `pg:"delegated"`
	Staking            string    `pg:"staking"`
	Unbonding          string    `pg:"unbonding"`
	Reward             string    `pg:"reward"`
	Address            string    `pg:"address"`
	Network            string    `pg:"network"`
	Significantdigits  int       `pg:"significantdigits"`
	Asset              string    `pg:"asset"`
}

func GetBalancesByDate(date time.Time) ([]BalanceAsset, error) {

	q := `
	SELECT b.*, a.Address, a.Network, a.SignificantDigits, a.Asset
	FROM Balances b
	JOIN Addresses a ON b.UUID = a.UUID
	WHERE DATE(b.Timestamp) > '%s'`

	q = fmt.Sprintf(q, date.Format("2006-01-02"))

	var results []BalanceAsset

	_, err := DB.Query(&results, q)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return nil, err
	}

	return results, nil
}
