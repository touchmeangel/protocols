package defillama

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Client struct {
	clients []*fasthttp.Client
}

func New(proxies []string) (*Client, error) {
	clients := make([]*fasthttp.Client, 0, len(proxies)+1)

	if len(proxies) == 0 {
		clients = append(clients, newClient(nil))
	} else {
		for _, proxy := range proxies {
			dialer := fasthttpproxy.FasthttpHTTPDialer(proxy)

			if len(proxy) >= 7 && proxy[:7] == "socks5:" {
				dialer = fasthttpproxy.FasthttpSocksDialer(proxy)
			}

			clients = append(clients, newClient(dialer))
		}
	}

	return &Client{
		clients: clients,
	}, nil
}

func newClient(dial fasthttp.DialFunc) *fasthttp.Client {
	return &fasthttp.Client{
		MaxConnsPerHost:     512,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
		MaxIdleConnDuration: 90 * time.Second,
		Dial:                dial,
	}
}

func (c *Client) randomClient() *fasthttp.Client {
	return c.clients[rand.IntN(len(c.clients))]
}

func (c *Client) get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	client := c.randomClient()

	if err := client.DoTimeout(req, resp, 10*time.Second); err != nil {
		return nil, err
	}

	body := append([]byte(nil), resp.Body()...)
	err := statusToError(body, resp.StatusCode())
	return body, err
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
