package models

import (
	"encoding/json"
	"time"
)

type Account struct {
	*AccountInfo
	Assets map[string]*AccountKDA `json:"assets"`
}

// AccountInfo holds (serializable) data about an account
type AccountInfo struct {
	Address       string        `json:"address,omitempty"`
	Nonce         uint64        `json:"nonce"`
	Name          string        `json:"name,omitempty"`
	RootHash      string        `json:"rootHash,omitempty"`
	Balance       int64         `json:"balance"`
	FrozenBalance int64         `json:"frozenBalance"`
	Allowance     int64         `json:"allowance"`
	Permissions   []Permissions `json:"permissions"`
	Timestamp     time.Duration `json:"timestamp"`
	Foundation    bool          `json:"foundation,omitempty"`
}

type PermissionKey struct {
	Address string `json:"address"`
	Weight  int64  `json:"weight"`
}

type Permissions struct {
	ID             int32           `json:"id"`
	Type           int32           `json:"type"`
	PermissionName string          `json:"permissionName"`
	Threshold      int64           `json:"Threshold"`
	Operations     string          `json:"operations"`
	Signers        []PermissionKey `json:"signers"`
}

func (acc *Account) String() string {
	result, err := json.MarshalIndent(acc, "", "\t")
	if err != nil {
		result = make([]byte, 0)
	}

	return string(result)
}
