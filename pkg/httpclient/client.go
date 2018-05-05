package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	*http.Client

	hdr http.Header
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Jar: new(Jar),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{},
			},
		},
	}
}

// WithProxy
func (c *Client) WithProxy(proxyURL *url.URL) *Client {
	if trans, ok := c.Client.Transport.(*http.Transport); ok {
		trans.Proxy = http.ProxyURL(proxyURL)
	}
	return c
}

// WithUserAgent
func (c *Client) WithUserAgent(ua string) *Client {
	if c.hdr == nil {
		c.hdr = http.Header{}
	}
	c.hdr.Set("User-Agent", ua)
	return c
}

// NewRequest
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.hdr {
		req.Header[k] = v
	}
	return req, nil
}

func (cli *Client) PostJSON(Url string, data interface{}) (*http.Response, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}

	req, err := cli.NewRequest("POST", Url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return cli.Do(req)
}

func (cli *Client) PostFile(Url string, formName string, f *os.File) (*http.Response, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	fw, err := w.CreateFormFile(formName, f.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}

	w.Close()
	req, err := cli.NewRequest("POST", Url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return cli.Do(req)
}
