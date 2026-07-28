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

	protocols, err := client.GetAllProtocols()
	if err != nil {
		log.Fatal(err)
	}

	filtered := defillama.Apply(protocols,
		defillama.ByMaxAudits(0),
		defillama.ByChain("Ethereum"),
		defillama.ByMinTVL(10_000),
		defillama.ByMaxTVL(300_000),
	)

	defillama.SortProtocols(filtered, defillama.SortByTVL, true)

	Print(filtered)
}

func Print(protocols []defillama.Protocol) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tNAME\tCATEGORY\tCHAINS\tTVL\t24H\t7D\tMCAP/TVL")
	for _, p := range protocols {
		category := "-"
		if p.Category != "" {
			category = p.Category
		}

		chains := strings.Join(p.Chains, ",")
		if len(chains) > 30 {
			chains = chains[:27] + "..."
		}

		ratio, ok := defillama.MCapToTVL(p)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t$%s\t%s\t%s\t%s\n",
			p.Slug,
			p.Name,
			category,
			chains,
			formatMoney(p.TVL),
			formatPct(p.Change1d),
			formatPct(p.Change7d),
			formatRatio(ratio, ok),
		)
	}
	fmt.Fprintf(w, "\nTotal: %d protocol(s)\n", len(protocols))
}

func formatPct(v float64) string {
	return fmt.Sprintf("%+.2f%%", v)
}

func formatRatio(v float64, ok bool) string {
	if !ok {
		return "-"
	}
	return fmt.Sprintf("%.2f", v)
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
