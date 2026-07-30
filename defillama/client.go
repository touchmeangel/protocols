package defillama

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const BaseURL = "https://api.llama.fi"
const DefaultClientTimeout = 10 * time.Second

type ProxyConfig struct {
	Address string
	Timeout time.Duration
}

type Client struct {
	baseURL       string
	client        *fasthttp.Client
	clientTimeout time.Duration
	proxy         string
}

func New() (*Client, error) {
	return &Client{baseURL: BaseURL, client: &fasthttp.Client{
		ReadTimeout:         DefaultClientTimeout,
		WriteTimeout:        DefaultClientTimeout,
		MaxConnsPerHost:     512,
		MaxIdleConnDuration: 90 * time.Second,
	}, clientTimeout: DefaultClientTimeout}, nil
}

// WithClientTimeout sets the per-request timeout (enforced via fasthttp's DoTimeout)
func (s *Client) WithClientTimeout(timeout time.Duration) *Client {
	s.clientTimeout = timeout
	s.client.ReadTimeout = timeout
	s.client.WriteTimeout = timeout
	if strings.HasPrefix(s.proxy, "http") {
		s.updateHTTPDialer(s.proxy, s.clientTimeout)
	}
	return s
}

// SetProxy
func (s *Client) SetProxy(proxyAddr string) error {
	if proxyAddr == "" {
		s.client.Dial = nil
		s.proxy = ""
		return nil
	}

	if strings.HasPrefix(proxyAddr, "http") {
		if _, err := url.Parse(proxyAddr); err != nil {
			return err
		}
		s.updateHTTPDialer(proxyAddr, s.clientTimeout)
		return nil
	}

	if strings.HasPrefix(proxyAddr, "socks") {
		if _, err := url.Parse(proxyAddr); err != nil {
			return fmt.Errorf("error creating socks proxy: %w", err)
		}
		s.client.Dial = fasthttpproxy.FasthttpSocksDialer(proxyAddr)
		s.proxy = proxyAddr
		return nil
	}

	return errors.New("only support http(s)/socks4/socks5 protocol")
}

func (s *Client) updateHTTPDialer(proxyAddr string, timeout time.Duration) {
	if timeout <= 0 {
		timeout = DefaultClientTimeout
	}
	s.client.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(proxyAddr, timeout)
	s.proxy = proxyAddr
}

func (c *Client) rawGet(endpoint string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + endpoint)
	req.Header.SetMethod(fasthttp.MethodGet)

	if err := c.client.DoTimeout(req, resp, c.clientTimeout); err != nil {
		return nil, fmt.Errorf("(timeout=%s): %w", c.clientTimeout, err)
	}

	body := append([]byte(nil), resp.Body()...)
	if err := statusToError(body, resp.StatusCode()); err != nil {
		return nil, err
	}
	return body, nil
}

func get[T any](c *Client, endpoint string) (T, error) {
	var out T
	raw, err := c.rawGet(endpoint)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("unmarshal %s: %w", endpoint, err)
	}
	return out, nil
}

func statusToError(body []byte, statusCode int) error {
	if statusCode != 200 {
		return fmt.Errorf("request failed with status %d: %s", statusCode, string(body))
	}
	return nil
}

func (c *Client) GetAllProtocols() ([]Protocol, error) {
	return get[[]Protocol](c, "/protocols")
}
