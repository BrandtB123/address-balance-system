package models

import "time"

type Bal struct {
	Address            string
	Timestamp          time.Time
	Balance            string
	Available          string
	DelegatableVesting string
	Delegated          string
	Staking            string
	Unbonding          string
	Reward             string
}

type Address struct {
	UUID              string
	Address           string
	Network           string
	SignificantDigits int
	Asset             string
}

type Test interface {
	NearValidatorAccountsResponse | NearBalanceResponse |
		AudiusResponse | OsmosisBalanceResponse |
		OsmosisDelegationResponse | OsmosisUnbondingResponse | OsmosisRewardResponse
}

type NearValidatorAccountsResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		BlockHash   string        `json:"block_hash"`
		BlockHeight int           `json:"block_height"`
		Logs        []interface{} `json:"logs"`
		Result      []int         `json:"result"`
	} `json:"result"`
	ID string `json:"id"`
}

type NearBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Amount        string `json:"amount"`
		BlockHash     string `json:"block_hash"`
		BlockHeight   int    `json:"block_height"`
		CodeHash      string `json:"code_hash"`
		Locked        string `json:"locked"`
		StoragePaidAt int    `json:"storage_paid_at"`
		StorageUsage  int    `json:"storage_usage"`
	} `json:"result"`
	ID string `json:"id"`
}

type AudiusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type OsmosisDelegationResponse struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type OsmosisRewardResponse struct {
	Rewards []struct {
		ValidatorAddress string `json:"validator_address"`
		Reward           []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"reward"`
	} `json:"rewards"`
	Total []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
}

type OsmosisUnbondingResponse struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type OsmosisBalanceResponse struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
