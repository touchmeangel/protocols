package defillama

type RootResponse struct {
	Props Props `json:"props"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	Chain     string     `json:"chain"`
	Metadata  Metadata   `json:"metadata"`
	Protocols []Protocol `json:"protocols"`
}

type Metadata struct {
	Name        string `json:"name"`
	Stablecoins bool   `json:"stablecoins"`
	Fees        bool   `json:"fees"`
	Dexs        bool   `json:"dexs"`
	Perps       bool   `json:"perps"`
	ID          string `json:"id"`
}

type Protocol struct {
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Logo           string     `json:"logo"`
	Category       *string    `json:"category"`
	Chains         []string   `json:"chains"`
	ChildProtocols []Protocol `json:"childProtocols,omitempty"`

	TVL         *TVL       `json:"tvl"`
	TVLChange   *TVLChange `json:"tvlChange"`
	MCap        *float64   `json:"mcap"`
	TokenPrice  *float64   `json:"tokenPrice"`
	MCapTVL     *float64   `json:"mcaptvl"`
	StrikeTVL   bool       `json:"strikeTvl"`
	DefiLlamaID string     `json:"defillamaId"`

	Fees           *Metric      `json:"fees"`
	Revenue        *Metric      `json:"revenue"`
	HoldersRevenue *Metric      `json:"holdersRevenue"`
	Emissions      *Metric      `json:"emissions"`
	Dexs           *DexMetrics  `json:"dexs"`
	Perps          *PerpMetrics `json:"perps,omitempty"`
}

type TVL struct {
	Default       *TVLValues `json:"default"`
	DoubleCounted *TVLValues `json:"doublecounted,omitempty"`
	Borrowed      *TVLValues `json:"borrowed,omitempty"`
	Staking       *TVLValues `json:"staking,omitempty"`
	Pool2         *TVLValues `json:"pool2,omitempty"`
}

type TVLValues struct {
	TVL          float64  `json:"tvl"`
	TVLPrevDay   *float64 `json:"tvlPrevDay"`
	TVLPrevWeek  *float64 `json:"tvlPrevWeek"`
	TVLPrevMonth *float64 `json:"tvlPrevMonth"`
}

type TVLChange struct {
	Change1D *float64 `json:"change1d"`
	Change7D *float64 `json:"change7d"`
	Change1M *float64 `json:"change1m"`
}

type Metric struct {
	Total24H         float64  `json:"total24h"`
	Total7D          float64  `json:"total7d"`
	Total30D         float64  `json:"total30d"`
	Total1Y          float64  `json:"total1y"`
	Annualized1Y     *float64 `json:"annualized1y"`
	MonthlyAverage1Y *float64 `json:"monthlyAverage1y"`
	TotalAllTime     float64  `json:"totalAllTime"`

	// Fees uses "pf", revenue uses "ps".
	PF *float64 `json:"pf,omitempty"`
	PS *float64 `json:"ps,omitempty"`
}

type DexMetrics struct {
	Total24H       float64 `json:"total24h"`
	Total7D        float64 `json:"total7d"`
	TotalAllTime   float64 `json:"totalAllTime"`
	Change7DOver7D float64 `json:"change_7dover7d"`
}

type PerpMetrics struct {
	Total24H       float64 `json:"total24h"`
	Total7D        float64 `json:"total7d"`
	Total30D       float64 `json:"total30d"`
	TotalAllTime   float64 `json:"totalAllTime"`
	Change7DOver7D float64 `json:"change_7dover7d"`
}
