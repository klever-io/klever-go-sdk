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
