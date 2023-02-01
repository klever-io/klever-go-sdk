package models

type KDAOptions struct {
	Name                 string
	Ticker               string
	Precision            int
	MaxSupply            float64
	InitialSupply        float64
	AddRolesMint         []string
	AddRolesSetITOPrices []string
	Roles                []*RolesInfo
	Properties           PropertiesInfo
	Attributes           AttributesInfo
	Staking              StakingInfo
	Royalties            RoyaltiesInfo

	Mime string
	Logo string
	URIs map[string]string
}

type AssetTriggerOptions struct {
	Amount               float64
	AddRolesMint         []string
	AddRolesSetITOPrices []string
	Staking              map[string]string
	Receiver             string
	Mime                 string
	Logo                 string
	URIs                 map[string]string
}

type ITOTriggerOptions struct {
	ReceiverAddress        string
	Status                 int32
	MaxAmount              float64
	PackInfo               map[string]PackInfoRequest
	DefaultLimitPerAddress int64
	WhitelistStatus        int32
	WhitelistInfo          map[string]WhitelistInfoRequest
	WhitelistStartTime     int64
	WhitelistEndTime       int64
	StartTime              int64
	EndTime                int64
}

type DepositType int32

const (
	FPRDeposit DepositType = iota
	KDAPoolDeposit
)

type DepositOptions struct {
	DepositType DepositType
	KDAID       string
	CurrencyID  string
	Amount      float64
}

type WithdrawType int32

const (
	StakingWithdraw WithdrawType = iota
	KDAPoolWithdraw
)

type WithdrawOptions struct {
	KDA          string
	WithdrawType WithdrawType
	Amount       float64
	CurrencyID   string
}
