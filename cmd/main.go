package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/touchmeangel/protocols/config"
	"github.com/touchmeangel/protocols/defillama"
)

func main() {
	cfg := config.Load()

	client, err := defillama.New(cfg.Proxies)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.GetAllProtocols()
	if err != nil {
		log.Fatal(err)
	}

	protocols := resp.Props.PageProps.Protocols

	filtered := defillama.Apply(protocols,
		defillama.ByChain("Ethereum"),
		defillama.ByMinTVL(10_000),
		defillama.ByMaxTVL(100_000),
	)

	defillama.SortProtocols(filtered, defillama.SortByTVL, true)
	top := defillama.Top(filtered, 20)

	Print(top)
}

func Print(protocols []defillama.Protocol) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NAME\tCATEGORY\tCHAINS\tTVL\t24H\t7D\tMCAP/TVL")
	for _, p := range protocols {
		category := "-"
		if p.Category != nil {
			category = *p.Category
		}

		tvl, _ := defillama.ProtocolTVL(p)

		chains := strings.Join(p.Chains, ",")
		if len(chains) > 30 {
			chains = chains[:27] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t$%s\t%s\t%s\t%s\n",
			p.Name,
			category,
			chains,
			formatMoney(tvl),
			formatPct(changeOrNil(p, defillama.SortByChange1D)),
			formatPct(changeOrNil(p, defillama.SortByChange7D)),
			formatRatio(p.MCapTVL),
		)
	}
	fmt.Fprintf(w, "\nTotal: %d protocol(s)\n", len(protocols))
}

func changeOrNil(p defillama.Protocol, field defillama.SortField) *float64 {
	if p.TVLChange == nil {
		return nil
	}
	switch field {
	case defillama.SortByChange1D:
		return p.TVLChange.Change1D
	case defillama.SortByChange7D:
		return p.TVLChange.Change7D
	}
	return nil
}

func formatPct(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%+.2f%%", *v)
}

func formatRatio(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.2f", *v)
}

func formatMoney(v float64) string {
	neg := v < 0
	if neg {
		v = -v
	}
	s := fmt.Sprintf("%.0f", v)

	n := len(s)
	if n <= 3 {
		if neg {
			return "-" + s
		}
		return s
	}

	var parts []string
	for n > 3 {
		parts = append([]string{s[n-3:]}, parts...)
		s = s[:n-3]
		n = len(s)
	}
	parts = append([]string{s}, parts...)

	out := strings.Join(parts, ",")
	if neg {
		out = "-" + out
	}
	return out
}
