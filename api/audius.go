package api

import (
	"fmt"
	"unit410/models"
)

type AudiusAPI struct{}

func (api *AudiusAPI) GetData() ([]*models.Bal, error) {
	var output []*models.Bal
	address := "0x46340b20830761efd32832a74d7169b29feb9758"
	for i := range [1000]int{} {
		fmt.Printf("%v: ", i)
		b, _ := getBalance2(address)
		fmt.Println(b)
		output = append(output, b)

	}
	return output, nil
}

func getBalance2(address string) (*models.Bal, error) {
	a, _ := HttpRequest[models.AudiusResponse]("GET", "https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=0x18aAA7115705e8be94bfFEBDE57Af9BFc265B998&address="+address+"&tag=latest&apikey=IQUJGKFGQSEUQMB8EKB8II6YTSY22ATSA3", nil)
	for a != nil && a.Status != "1" {
		a, _ = HttpRequest[models.AudiusResponse]("GET", "https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=0x18aAA7115705e8be94bfFEBDE57Af9BFc265B998&address="+address+"&tag=latest&apikey=IQUJGKFGQSEUQMB8EKB8II6YTSY22ATSA3", nil)
	}
	balance := models.Bal{
		Address: address,
		Balance: a.Result,
	}

	return &balance, nil
}
