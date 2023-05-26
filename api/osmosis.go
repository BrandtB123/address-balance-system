package api

import (
	"fmt"
	"math/big"
	"unit410/db"
	"unit410/models"
)

type OsmosisAPI struct{}

var osmosisAddresses = []string{
	"osmo15ujzuqf7594x7wqssgqv0tcrmjccyhwy9wuyw0",
	"osmo1x0l3jhq4k0uv0qwatvzk97ktmk5uewclnzsfpl",
	"osmo1npfl4vmmmf4yqhcemz95mvqujgdnlhrlxfzhlhz2gru8g2t749xqr9zm5e",
	"osmo1exeej95hmmzavtzsjnrgpkhrj0cvp72lvl9qf5",
	"osmo1qvsus5qg8yhre7k2c78xkkw4nvqqgev7dvylaj",
}

func (api *OsmosisAPI) GetData() error {
	for _, address := range osmosisAddresses {
		err := getOsmosisBalance(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertListToBigInt(inputs []string) string {
	sum := big.NewInt(0)
	for _, numStr := range inputs {
		num := new(big.Int)
		num.SetString(numStr, 10)
		sum.Add(sum, num)
	}
	return sum.String()
}

func getOsmosisBalance(address string) error {
	osmosisBalanceResponse, err := HttpRequest[models.OsmosisBalanceResponse]("GET", "https://lcd.osmosis.zone/cosmos/bank/v1beta1/balances/"+address, nil)
	if err != nil {
		return err
	}
	var balance models.Bal
	for _, bal := range osmosisBalanceResponse.Balances {
		asset := bal.Denom
		balance = models.Bal{
			Address: address,
			Balance: bal.Amount,
		}
		err = db.AddAddress(models.Address{
			Address:           address,
			Network:           "Osmosis",
			SignificantDigits: 18,
			Asset:             asset,
		})
		if asset != "uosmo" {
			err = db.AddBalance(asset, balance)
			if err != nil {
				return err
			}
		}
		fmt.Println(err)
	}
	delegationResponse, err := HttpRequest[models.OsmosisDelegationResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegations/"+address, nil)
	if err != nil {
		return err
	}
	var delegatedBalances []string
	for _, delegated := range delegationResponse.DelegationResponses {
		if delegated.Delegation.DelegatorAddress == address {
			delegatedBalances = append(delegatedBalances, delegated.Balance.Amount)
		}
	}
	balance.Delegated = convertListToBigInt(delegatedBalances)

	unbondingResponse, err := HttpRequest[models.OsmosisUnbondingResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegators/"+address+"/unbonding_delegations", nil)
	if err != nil {
		return err
	}
	var unbondedBalances []string
	for _, unbonded := range unbondingResponse.UnbondingResponses {
		if unbonded.DelegatorAddress == address {
			for _, entry := range unbonded.Entries {
				unbondedBalances = append(unbondedBalances, entry.Balance)
			}
		}
	}
	balance.Unbonding = convertListToBigInt(unbondedBalances)

	rewardResponse, err := HttpRequest[models.OsmosisRewardResponse]("GET", "https://lcd.osmosis.zone/cosmos/distribution/v1beta1/delegators/"+address+"/rewards", nil)
	if err != nil {
		return err
	}
	balance.Reward = rewardResponse.Total[0].Amount
	db.AddBalance("uosmo", balance)
	return nil
}
