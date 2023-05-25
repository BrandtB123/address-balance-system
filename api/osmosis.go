package api

import (
	"fmt"
	"unit410/models"
)

type OsmosisAPI struct{}

func (api *OsmosisAPI) GetData() ([]*models.Bal, error) {
	address := "osmo1npfl4vmmmf4yqhcemz95mvqujgdnlhrlxfzhlhz2gru8g2t749xqr9zm5e"
	getBalance(address)
	return nil, nil
}

func getBalance(address string) {
	//var balance models.Bal
	a, err := HttpRequest[models.OsmosisBalanceResponse]("GET", "https://lcd.osmosis.zone/cosmos/bank/v1beta1/balances/"+address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	x, err := HttpRequest[models.OsmosisDelegationResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegations/"+address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	b, err := HttpRequest[models.OsmosisUnbondingResponse]("GET", "https://lcd.osmosis.zone/cosmos/staking/v1beta1/delegators/"+address+"/unbonding_delegations", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	c, err := HttpRequest[models.OsmosisRewardResponse]("GET", "https://lcd.osmosis.zone/cosmos/distribution/v1beta1/delegators/"+address+"/rewards", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a)
	fmt.Println(x)
	fmt.Println(b)
	fmt.Println(c)
}
