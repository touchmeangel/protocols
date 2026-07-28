package defillama

import "sort"

type SortField int

const (
	SortByTVL SortField = iota
	SortByChange1D
	SortByChange7D
	SortByMCap
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
		a, b := sortValue(protocols[i], field), sortValue(protocols[j], field)
		if desc {
			return a > b
		}
		return a < b
	})
}

func sortValue(p Protocol, field SortField) float64 {
	switch field {
	case SortByTVL:
		return p.TVL
	case SortByChange1D:
		return p.Change1d
	case SortByChange7D:
		return p.Change7d
	case SortByMCap:
		if p.MCap != nil {
			return *p.MCap
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
