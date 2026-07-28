package defillama

import "strings"

type Filter func(Protocol) bool

func Apply(protocols []Protocol, filters ...Filter) []Protocol {
	if len(filters) == 0 {
		return protocols
	}
	out := make([]Protocol, 0, len(protocols))
	for _, p := range protocols {
		if matchesAll(p, filters) {
			out = append(out, p)
		}
	}
	return out
}

func matchesAll(p Protocol, filters []Filter) bool {
	for _, f := range filters {
		if !f(p) {
			return false
		}
	}
	return true
}

func Any(filters ...Filter) Filter {
	return func(p Protocol) bool {
		for _, f := range filters {
			if f(p) {
				return true
			}
		}
		return false
	}
}

func Not(f Filter) Filter {
	return func(p Protocol) bool {
		return !f(p)
	}
}

func ByChain(chain string) Filter {
	chain = strings.ToLower(chain)
	return func(p Protocol) bool {
		for _, c := range p.Chains {
			if strings.ToLower(c) == chain {
				return true
			}
		}
		return false
	}
}

func ByCategory(category string) Filter {
	category = strings.ToLower(category)
	return func(p Protocol) bool {
		if p.Category == nil {
			return false
		}
		return strings.ToLower(*p.Category) == category
	}
}

func ByNameContains(substr string) Filter {
	substr = strings.ToLower(substr)
	return func(p Protocol) bool {
		return strings.Contains(strings.ToLower(p.Name), substr)
	}
}

func ProtocolTVL(p Protocol) (float64, bool) {
	if p.TVL == nil || p.TVL.Default == nil {
		return 0, false
	}
	return p.TVL.Default.TVL, true
}

func ByMinTVL(min float64) Filter {
	return func(p Protocol) bool {
		tvl, ok := ProtocolTVL(p)
		return ok && tvl >= min
	}
}

func ByMaxTVL(max float64) Filter {
	return func(p Protocol) bool {
		tvl, ok := ProtocolTVL(p)
		return ok && tvl <= max
	}
}

func ByTVLRange(min, max float64) Filter {
	return func(p Protocol) bool {
		tvl, ok := ProtocolTVL(p)
		return ok && tvl >= min && tvl <= max
	}
}

func ByMinChange1D(min float64) Filter {
	return func(p Protocol) bool {
		if p.TVLChange == nil || p.TVLChange.Change1D == nil {
			return false
		}
		return *p.TVLChange.Change1D >= min
	}
}

func ByMaxChange1D(max float64) Filter {
	return func(p Protocol) bool {
		if p.TVLChange == nil || p.TVLChange.Change1D == nil {
			return false
		}
		return *p.TVLChange.Change1D <= max
	}
}

func ByMinChange7D(min float64) Filter {
	return func(p Protocol) bool {
		if p.TVLChange == nil || p.TVLChange.Change7D == nil {
			return false
		}
		return *p.TVLChange.Change7D >= min
	}
}

func ByMaxChange7D(max float64) Filter {
	return func(p Protocol) bool {
		if p.TVLChange == nil || p.TVLChange.Change7D == nil {
			return false
		}
		return *p.TVLChange.Change7D <= max
	}
}

func ByMinFees24H(min float64) Filter {
	return func(p Protocol) bool {
		return p.Fees != nil && p.Fees.Total24H >= min
	}
}

func ByMinRevenue24H(min float64) Filter {
	return func(p Protocol) bool {
		return p.Revenue != nil && p.Revenue.Total24H >= min
	}
}

func ByMinMCap(min float64) Filter {
	return func(p Protocol) bool {
		return p.MCap != nil && *p.MCap >= min
	}
}

func ByMaxMCapTVLRatio(max float64) Filter {
	return func(p Protocol) bool {
		return p.MCapTVL != nil && *p.MCapTVL > 0 && *p.MCapTVL <= max
	}
}
