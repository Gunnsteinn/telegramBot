package domain

// Sponsor (Sponsor info) response is exported, it models the data we receive.
type Sponsor struct {
	Name          string `bson:"name,omitempty" json:"name,omitempty"`
	LastName      string `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Email         string `bson:"email,omitempty" json:"email,omitempty"`
	WalletAddress string `bson:"wallet_address,omitempty" json:"wallet_address,omitempty"`
	Teams         []struct {
		TeamName      string  `bson:"team_name,omitempty" json:"team_name,omitempty"`
		WalletAddress string  `bson:"wallet_address,omitempty" json:"wallet_address,omitempty"`
		PoolPercent   float64 `bson:"pool_percent,omitempty" json:"pool_percent,omitempty"`
		Adventurer    struct {
			Name              string  `bson:"name,omitempty" json:"name,omitempty"`
			LastName          string  `bson:"last_name,omitempty" json:"last_name,omitempty"`
			WalletAddress     string  `bson:"wallet_address,omitempty" json:"wallet_address,omitempty"`
			ProfitSlp         int     `bson:"profit_slp,omitempty" json:"profit_slp,omitempty"`
			Performance       float64 `bson:"performance,omitempty" json:"performance,omitempty"`
			LastClaimedItemAt string  `bson:"last_claimed_item_at,omitempty" json:"last_claimed_item_at,omitempty"`
			PlayedDays        float64 `bson:"played_days,omitempty" json:"played_days,omitempty"`
		} `bson:"Adventurer,omitempty" json:"Adventurer,omitempty"`
	} `bson:"teams,omitempty" json:"teams,omitempty"`
}
