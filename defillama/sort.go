package defillama

import "sort"

type SortField int

const (
	SortByTVL SortField = iota
	SortByChange1D
	SortByChange7D
	SortByMCap
	SortByFees24H
	SortByRevenue24H
	SortByName
)

func SortProtocols(protocols []Protocol, field SortField, desc bool) {
	sort.SliceStable(protocols, func(i, j int) bool {
		if field == SortByName {
			if desc {
				return protocols[i].Name > protocols[j].Name
			}
			return protocols[i].Name < protocols[j].Name
		}

		a := sortValue(protocols[i], field)
		b := sortValue(protocols[j], field)
		if desc {
			return a > b
		}
		return a < b
	})
}

func sortValue(p Protocol, field SortField) float64 {
	switch field {
	case SortByTVL:
		v, _ := ProtocolTVL(p)
		return v
	case SortByChange1D:
		if p.TVLChange != nil && p.TVLChange.Change1D != nil {
			return *p.TVLChange.Change1D
		}
	case SortByChange7D:
		if p.TVLChange != nil && p.TVLChange.Change7D != nil {
			return *p.TVLChange.Change7D
		}
	case SortByMCap:
		if p.MCap != nil {
			return *p.MCap
		}
	case SortByFees24H:
		if p.Fees != nil {
			return p.Fees.Total24H
		}
	case SortByRevenue24H:
		if p.Revenue != nil {
			return p.Revenue.Total24H
		}
	}
	return 0
}

func Top(protocols []Protocol, n int) []Protocol {
	if n >= len(protocols) {
		return protocols
	}
	return protocols[:n]
}
