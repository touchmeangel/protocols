package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

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

	client, err := defillama.New()
	if err != nil {
		log.Fatal(err)
	}

	proxy := cfg.Proxies[rand.IntN(len(cfg.Proxies))]
	client.WithClientTimeout(proxy.Timeout)
	client.SetProxy(proxy.Address)
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

	followerFilters := []twitter.FollowerFilter{
		twitter.MinFollowers(100),
		twitter.MaxFollowers(10000),
	}

	results := fetchByFollowers(filtered, cfg.Proxies, followerFilters)

	Print(results)
}

const concurrentFetches = 5

type followerOutcome struct {
	protocol defillama.Protocol
	matched  bool
}

const profileAttempts = 3

func fetchProfileWithRetry(client *twitter_api.Client, proxies []defillama.ProxyConfig, username string) (twitter_api.Profile, error) {
	var lastErr error

	for attempt := 0; attempt < profileAttempts; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt)*time.Second + time.Duration(rand.IntN(500))*time.Millisecond
			time.Sleep(backoff)

			proxy := proxies[rand.IntN(len(proxies))]
			client.WithClientTimeout(proxy.Timeout)
			if err := client.SetProxy(proxy.Address); err != nil {
				lastErr = err
				continue
			}
			client.ResetGuestSession()
		}

		profile, err := client.GetProfile(username)
		if err == nil {
			return profile, nil
		}
		lastErr = err

		var apiErr *twitter_api.APIError
		if errors.As(err, &apiErr) && (apiErr.StatusCode == http.StatusTooManyRequests || apiErr.StatusCode == http.StatusForbidden) {
			continue
		}
		return twitter_api.Profile{}, err
	}

	return twitter_api.Profile{}, fmt.Errorf("gave up after %d attempts: %w", profileAttempts, lastErr)
}

func fetchByFollowers(filtered []defillama.Protocol, proxies []defillama.ProxyConfig, followerFilters []twitter.FollowerFilter) []defillama.Protocol {
	type job struct {
		index    int
		protocol defillama.Protocol
	}

	jobs := make(chan job)
	outcomes := make([]followerOutcome, len(filtered))

	var wg sync.WaitGroup
	for i := 0; i < concurrentFetches; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := twitter_api.New()

			for j := range jobs {
				p := j.protocol
				proxy := proxies[rand.IntN(len(proxies))]
				err := client.SetProxy(proxy.Address)
				if err != nil {
					log.Printf("%s proxy fails with: %v", proxy.Address, err)
					continue
				}
				client.WithClientTimeout(proxy.Timeout)

				profile, err := fetchProfileWithRetry(client, proxies, p.Twitter)
				if err != nil {
					switch {
					case strings.Contains(err.Error(), "rest_id not found"):
						log.Printf("twitter handle %q: account not found / suspended", p.Twitter)
					default:
						log.Printf("twitter handle %q: fetch failed: %v", p.Twitter, err)
					}
					continue
				}

				if twitter.MatchesAllFollowers(profile.FollowersCount, followerFilters) {
					fmt.Printf("%s: %d followers\n", p.Twitter, profile.FollowersCount)
					outcomes[j.index] = followerOutcome{protocol: p, matched: true}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for i, p := range filtered {
			if p.Twitter == "" {
				continue
			}
			jobs <- job{index: i, protocol: p}
		}
	}()

	wg.Wait()

	results := make([]defillama.Protocol, 0, len(filtered))
	for _, o := range outcomes {
		if o.matched {
			results = append(results, o.protocol)
		}
	}
	return results
}

func Print(protocols []defillama.Protocol) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	_, _ = fmt.Fprintln(w, "ID\tNAME\tCATEGORY\tCHAINS\tTVL\t24H\t7D\tMCAP/TVL")
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

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t$%s\t%s\t%s\t%s\n",
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
	_, _ = fmt.Fprintf(w, "\nTotal: %d protocol(s)\n", len(protocols))
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
