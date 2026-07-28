package defillama

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"golang.org/x/net/http/httpproxy"
)

type ProxyConfig struct {
	Address string
	Timeout time.Duration
}

type proxyClient struct {
	client  *fasthttp.Client
	label   string
	timeout time.Duration
}

type Client struct {
	clients []proxyClient
}

func New(proxies []ProxyConfig) (*Client, error) {
	clients := make([]proxyClient, 0, len(proxies))

	for _, p := range proxies {
		label := p.Address
		var dial fasthttp.DialFunc

		switch {
		case p.Address == "":
			label = "direct"
		case len(p.Address) >= 7 && p.Address[:7] == "socks5:":
			dial = socksDialerWithTimeout(p.Address, p.Timeout)
		default:
			dial = fasthttpproxy.FasthttpHTTPDialerTimeout(p.Address, p.Timeout)
		}

		clients = append(clients, proxyClient{
			client:  newClient(dial, p.Timeout),
			label:   label,
			timeout: p.Timeout,
		})
	}

	if len(clients) == 0 {
		return nil, errors.New("no proxy/client entries configured")
	}

	return &Client{clients: clients}, nil
}

func newClient(dial fasthttp.DialFunc, timeout time.Duration) *fasthttp.Client {
	return &fasthttp.Client{
		Dial:                dial,
		ReadTimeout:         timeout,
		WriteTimeout:        timeout,
		MaxConnsPerHost:     512,
		MaxIdleConnDuration: 90 * time.Second,
	}
}

func socksDialerWithTimeout(proxyAddr string, timeout time.Duration) fasthttp.DialFunc {
	d := fasthttpproxy.Dialer{
		Config:         httpproxy.Config{HTTPProxy: proxyAddr, HTTPSProxy: proxyAddr},
		Timeout:        timeout,
		ConnectTimeout: timeout,
	}
	dialFunc, _ := d.GetDialFunc(false)
	return dialFunc
}

func (c *Client) randomClient() proxyClient {
	return c.clients[rand.IntN(len(c.clients))]
}

func (c *Client) get(url string) ([]byte, error) {
	pc := c.randomClient()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	if err := pc.client.DoTimeout(req, resp, pc.timeout); err != nil {
		return nil, fmt.Errorf("via %s (timeout=%s): %w", pc.label, pc.timeout, err)
	}

	body := append([]byte(nil), resp.Body()...)
	err := statusToError(body, resp.StatusCode())
	if err != nil {
		return body, err
	}

	raw, err := Extract(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func statusToError(body []byte, statusCode int) error {
	if statusCode != 200 {
		return fmt.Errorf("request failed with status %d: %s", statusCode, string(body))
	}
	return nil
}

func (c *Client) GetAllProtocols() error {
	body, err := c.get("https://defillama.com/")
	if err != nil {
		return err
	}

	fmt.Printf("resp: %s", body)
	return nil
}
