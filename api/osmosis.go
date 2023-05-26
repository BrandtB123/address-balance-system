package api

import (
	"fmt"
	"unit410/db"
	"unit410/models"
)

type OsmosisAPI struct{}

func (api *OsmosisAPI) GetData() ([]*models.Bal, error) {
	address := "osmo1gy0nyn2hscxxayj2pdyu8axmfvv75nnvd08ajh"
	for i := range [10]int{} {
		fmt.Printf("%v: ", i)
		_, err := getBalance(address)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func getBalance(address string) (*models.Bal, error) {
	//var balance models.Bal
	a, err := HttpRequest[models.OsmosisBalanceResponse]("GET", "https://lcd.osmosis.zone/cosmos/bank/v1beta1/balances/"+address, nil)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err)
		err = db.AddBalance(b.Denom, balance)
		fmt.Println(err)
	}
	// have to do something here about how the balances returns a list for all assets but the next 3 methods only do 1 asset
	return nil, nil
	// x, err := HttpRequest[models.OsmosisDelegationResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegations/"+address, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// b, err := HttpRequest[models.OsmosisUnbondingResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegators/"+address+"/unbonding_delegations", nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// c, err := HttpRequest[models.OsmosisRewardResponse]("GET", "https://lcd.osmosis.zone/cosmos/distribution/v1beta1/delegators/"+address+"/rewards", nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}
