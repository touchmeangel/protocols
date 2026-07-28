package defillama

import "encoding/json"

type Protocol struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Address     string   `json:"address"`
	Symbol      string   `json:"symbol"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Chain       string   `json:"chain"`
	Logo        string   `json:"logo"`
	Category    string   `json:"category"`
	Chains      []string `json:"chains"`
	Module      string   `json:"module"`
	Twitter     string   `json:"twitter"`
	ListedAt    int64    `json:"listedAt"`
	Audits      string   `json:"audits"`
	GeckoID     string   `json:"gecko_id"`
	CmcID       string   `json:"cmcId"`
	Methodology string   `json:"methodology,omitempty"`
	TvlCodePath string   `json:"tvlCodePath,omitempty"`

	TVL       float64            `json:"tvl"`
	ChainTvls map[string]float64 `json:"chainTvls"`
	MCap      *float64           `json:"mcap"`

	Change1h float64 `json:"change_1h"`
	Change1d float64 `json:"change_1d"`
	Change7d float64 `json:"change_7d"`

	AuditLinks           []string `json:"audit_links,omitempty"`
	OpenSource           bool     `json:"openSource,omitempty"`
	GitHub               []string `json:"github,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	PreviousNames        []string `json:"previousNames,omitempty"`
	ParentProtocol       string   `json:"parentProtocol,omitempty"`
	ParentProtocolSlug   string   `json:"parentProtocolSlug,omitempty"`
	Deprecated           bool     `json:"deprecated,omitempty"`
	ReferralUrl          string   `json:"referralUrl,omitempty"`
	AssetToken           string   `json:"assetToken,omitempty"`
	Staking              *float64 `json:"staking,omitempty"`
	Borrowed             *float64 `json:"borrowed,omitempty"`
	Pool2                *float64 `json:"pool2,omitempty"`
	MisrepresentedTokens bool     `json:"misrepresentedTokens,omitempty"`

	Dimensions       json.RawMessage `json:"dimensions,omitempty"`
	Hallmarks        json.RawMessage `json:"hallmarks,omitempty"`
	OraclesBreakdown json.RawMessage `json:"oraclesBreakdown,omitempty"`
	TokenBreakdowns  json.RawMessage `json:"tokenBreakdowns,omitempty"`
}
