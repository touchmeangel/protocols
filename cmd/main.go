package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/touchmeangel/protocols/config"
	"github.com/touchmeangel/protocols/defillama"
	"github.com/touchmeangel/protocols/twitter"
	twitter_api "github.com/touchmeangel/twitter_api"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

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

	twitterClient := twitter_api.New()

	followerFilters := []twitter.FollowerFilter{
		twitter.MinFollowers(100),
		twitter.MaxFollowers(10000),
	}

	var results []defillama.Protocol

	for _, p := range filtered {
		if p.Twitter == "" {
			continue
		}

		acc := cfg.Accounts[rand.IntN(len(cfg.Accounts))]
		twitterClient.SetAuthToken(acc.Token)

		profile, err := twitterClient.GetProfile(p.Twitter)
		if err != nil {
			if strings.Contains(err.Error(), "rest_id not found") {
				log.Printf("twitter handle %q: account not found / suspended", p.Twitter)
				continue
			}
			log.Printf("twitter handle %q: fetch failed using account %q: %v", p.Twitter, acc.Label, err)
			continue
		}

		if !twitter.MatchesAllFollowers(profile.FollowersCount, followerFilters) {
			continue
		}

		fmt.Printf("%s: %d followers\n", p.Twitter, profile.FollowersCount)
		results = append(results, p)
	}

	Print(results)
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
