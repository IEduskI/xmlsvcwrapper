package xmlsvcwrapper

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	httpClient *http.Client
}

// New function creates a new client http instance
func New() *Client {
	return NewClient(&http.Client{})
}

func NewClient(c *http.Client) *Client {
	if c.Transport == nil {
		c.Transport = createTransport(nil)
	}
	return &Client{
		httpClient: c,
	}
}

// R function creates a new request instance
func (c *Client) R() *Request {
	return &Request{
		client: c,
		Header: http.Header{},
	}
}

// SetTimeOut method sets timeout for request raised from client.
//
//	client.SetTimeout(time.Duration(1 * time.Second))
func (c *Client) SetTimeOut(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func createTransport(httpTransport *http.Transport) *http.Transport {

	if httpTransport != nil {
		return httpTransport
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

}