package defillama

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	client *fasthttp.Client
}

func New() (*Client, error) {
	client := &fasthttp.Client{
		MaxConnsPerHost:     512,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
		MaxIdleConnDuration: 90 * time.Second,
		// Dial: fasthttpproxy.FasthttpSocksDialer("socks5://myuser:mypass@127.0.0.1:1080"),
	}

	return &Client{client: client}, nil
}

// https://defillama.com/protocol/surf-lending
func (c *Client) get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	if err := c.client.DoTimeout(req, resp, 10*time.Second); err != nil {
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
