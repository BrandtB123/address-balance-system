package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"unit410/db"
	"unit410/models"
)

type NearAPI struct {
}

func (api *NearAPI) GetData() ([]*models.Bal, error) {
	stakedBalances, _ := obtainValidatorAddresses()
	for i, sb := range stakedBalances {
		b, err := getAccountBalance(sb.Address)
		db.AddAddress(models.Address{
			Address:           sb.Address,
			SignificantDigits: 112,
			Network:           "NEAR",
			Asset:             "NEAR",
		})
		fmt.Print(i)
		fmt.Println(" ", b)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		stakedBalances[i].Balance = b.Result.Amount
		db.AddBalance("NEAR", *stakedBalances[i])
	}
	fmt.Println(len(stakedBalances))

	return stakedBalances, nil
}

type Account struct {
	AccountID       string `json:"account_id"`
	UnstakedBalance string `json:"unstaked_balance"`
	StakedBalance   string `json:"staked_balance"`
	CanWithdraw     bool   `json:"can_withdraw"`
}

func parseAccountsJSON(jsonStr string) ([]Account, error) {
	var accounts []Account
	err := json.Unmarshal([]byte(jsonStr), &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func asciiCodesToString(asciiCodes []int) string {
	var runes []rune

	for _, code := range asciiCodes {
		runes = append(runes, rune(code))
	}

	return string(runes)
}

func JSONToBase64(jsonData map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str, nil
}

func obtainValidatorAddresses() ([]*models.Bal, error) {
	fromIndex := 0
	limit := 50
	increment := 50
	toIndex := 100

	jsonData := map[string]interface{}{
		"from_index": fromIndex,
		"limit":      limit,
	}
	var balances []*models.Bal

	for fromIndex+increment <= toIndex {
		fromIndex += increment

		jsonData["from_index"] = fromIndex
		jsonData["limit"] = limit

		base64Str, err := JSONToBase64(jsonData)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		requestBody2 := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      "id",
			"method":  "query",
			"params": map[string]interface{}{
				"request_type": "call_function",
				"finality":     "final",
				"account_id":   "astro-stakers.poolv1.near",
				"method_name":  "get_accounts",
				"args_base64":  base64Str,
			},
		}

		y, e := HttpRequest[models.NearValidatorAccountsResponse]("POST", "https://rpc.mainnet.near.org", requestBody2)
		if e != nil {
			fmt.Println(e.Error())
		}
		a := asciiCodesToString(y.Result.Result)
		b, _ := parseAccountsJSON(a)

		for _, account := range b {
			temp := models.Bal{
				Address: account.AccountID,
				Staking: account.StakedBalance,
			}
			balances = append(balances, &temp)
		}
	}
	return balances, nil
}

func getAccountBalance(address string) (*models.NearBalanceResponse, error) {
	// Construct the NEAR JSON-RPC API request URL
	rpcURL := "https://rpc.mainnet.near.org"

	// Prepare the JSON-RPC request payloa
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "id",
		"method":  "query",
		"params": map[string]interface{}{
			"request_type": "view_account",
			"finality":     "final",
			"account_id":   address,
		},
	}

	// Convert the request payload to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Send the JSON-RPC request
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	// Parse the JSON response
	var response *models.NearBalanceResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
