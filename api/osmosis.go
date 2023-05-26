package api

import (
	"fmt"
	"math/big"
	"unit410/db"
	"unit410/models"
)

type OsmosisAPI struct{}

func (api *OsmosisAPI) GetData() error {
	address := "osmo1gy0nyn2hscxxayj2pdyu8axmfvv75nnvd08ajh"
	for range [10]int{} {
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
	a, err := HttpRequest[models.OsmosisBalanceResponse]("GET", "https://lcd.osmosis.zone/cosmos/bank/v1beta1/balances/"+address, nil)
	if err != nil {
		return err
	}
	var balance models.Bal
	for _, b := range a.Balances {
		balance = models.Bal{
			Address: address,
			Balance: b.Amount,
		}
		err = db.AddAddress(models.Address{
			Address:           address,
			Network:           "Osmosis",
			SignificantDigits: 10,
			Asset:             b.Denom,
		})
		if b.Denom != "uosmo" {
			err = db.AddBalance(b.Denom, balance)
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
