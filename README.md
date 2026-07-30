# protogod

> Web3 protocol data scraper that aggregates DefiLlama metrics and social signals from X (Twitter).

Created by **touchmeangel**

## Quick Start

```bash
git clone https://github.com/touchmeangel/protogod.git
cd protogod

go run ./cmd
```

### Features

* Collect protocol metrics from DefiLlama
* Scrape protocol followers count from X (Twitter)
* Normalize collected data
* Easy to extend with additional data sources

### Configuration

Create a `config.json` with proxies you'd like to use if any

```json
{
    "proxies": [
        "socks5://user:pass@163.40.2.10"
    ]
}
```

### Roadmap

* [x] DefiLlama scraper
* [x] X (Twitter) scraper
* [x] Docker support
* [ ] Scheduled scraping
* [ ] Telegram bot integration
* [ ] Real-time Telegram alerts for protocol updates

### Example Output

```bash
ID                     NAME                   CATEGORY                   CHAINS                          TVL       24H        7D         MCAP/TVL  FOLLOWERS
swellchain-bridge      Swellchain Bridge      Canonical Bridge           Ethereum                        $290,880  +0.74%     -1.20%     -         510
union-protocol         Union Protocol         Uncollateralized Lending   Optimism,Ethereum,Base,Arbi...  $266,394  -0.07%     -0.06%     -         5,753
sectorone-dlmm         SectorOne DLMM         Dexs                       Robinhood Chain,Base,MegaET...  $251,005  -0.08%     -2.60%     -         1,156
solid-yield            Solid Yield            Yield                      Fuse,Ethereum                   $249,348  +0.10%     -0.17%     -         3,692
mim-swap               MIM Swap               Dexs                       Arbitrum,Blast,Kava,Nibiru,...  $242,227  +39.96%    -10.82%    -         545
monolith-market        Monolith Market        CDP                        Ethereum                        $209,414  +0.00%     -6.89%     -         248
steroids               Steroids               Yield                      Ethereum                        $208,808  +0.17%     -2.02%     -         104
unipower               UniPower               Yield                      Ethereum,Polygon                $204,184  +0.84%     -0.81%     -         1,063
surf-liquid            Surf Liquid            Yield                      Base,Polygon,Ethereum,Arbitrum  $204,074  +0.07%     -0.64%     0.54      2,995
shimmerbridge          ShimmerBridge          Bridge                     Binance,Ethereum,Avalanche,...  $203,263  +0.04%     -1.20%     -         726
nsure                  Nsure                  Insurance                  Ethereum                        $194,948  -0.28%     -1.05%     0.10      5,914
```

### Collaborators

* **[touchmeangel](https://github.com/touchmeangel)** — Creator
