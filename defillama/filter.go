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
	return func(p Protocol) bool { return !f(p) }
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
	return func(p Protocol) bool { return strings.ToLower(p.Category) == category }
}

func ByNameContains(substr string) Filter {
	substr = strings.ToLower(substr)
	return func(p Protocol) bool { return strings.Contains(strings.ToLower(p.Name), substr) }
}

func ByMinTVL(min float64) Filter { return func(p Protocol) bool { return p.TVL >= min } }
func ByMaxTVL(max float64) Filter { return func(p Protocol) bool { return p.TVL <= max } }
func ByTVLRange(min, max float64) Filter {
	return func(p Protocol) bool { return p.TVL >= min && p.TVL <= max }
}

func ByMinChange1D(min float64) Filter { return func(p Protocol) bool { return p.Change1d >= min } }
func ByMaxChange1D(max float64) Filter { return func(p Protocol) bool { return p.Change1d <= max } }
func ByMinChange7D(min float64) Filter { return func(p Protocol) bool { return p.Change7d >= min } }
func ByMaxChange7D(max float64) Filter { return func(p Protocol) bool { return p.Change7d <= max } }

func ByMinMCap(min float64) Filter {
	return func(p Protocol) bool { return p.MCap != nil && *p.MCap >= min }
}

func MCapToTVL(p Protocol) (float64, bool) {
	if p.MCap == nil || *p.MCap == 0 || p.TVL == 0 {
		return 0, false
	}
	return *p.MCap / p.TVL, true
}

func ByMaxMCapTVLRatio(max float64) Filter {
	return func(p Protocol) bool {
		ratio, ok := MCapToTVL(p)
		return ok && ratio <= max
	}
}
