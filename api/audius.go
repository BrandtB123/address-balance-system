package api

import (
	"unit410/db"
	"unit410/models"
)

type AudiusAPI struct{}

var audius = "0x18aAA7115705e8be94bfFEBDE57Af9BFc265B998"
var apiKey = "IQUJGKFGQSEUQMB8EKB8II6YTSY22ATSA3"

var audiusAddresses = []string{
	"0x28c6c06298d514db089934071355e5743bf21d60",
	"0x46340b20830761efd32832a74d7169b29feb9758"}

func (api *AudiusAPI) GetData() error {
	for _, addr := range audiusAddresses {
		balance, _ := getAudiusBalance(addr)
		db.AddAddress(models.Address{
			Address:           addr,
			Network:           "eth",
			SignificantDigits: 18,
			Asset:             audius,
		})
		db.AddBalance(audius, *balance)

	}
	return nil
}

func getAudiusBalance(address string) (*models.Bal, error) {
	audiusResponse, _ := HttpRequest[models.AudiusResponse]("GET", "https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress="+audius+"&address="+address+"&tag=latest&apikey="+apiKey, nil)
	for audiusResponse != nil && audiusResponse.Status != "1" {
		audiusResponse, _ = HttpRequest[models.AudiusResponse]("GET", "https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress="+audius+"&address="+address+"&tag=latest&apikey="+apiKey, nil)
	}
	balance := models.Bal{
		Address: address,
		Balance: audiusResponse.Result,
	}

	return &balance, nil
}
