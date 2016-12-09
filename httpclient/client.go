package httpclient

import (
	"net/http"
	"net/url"
)

type CookieMap map[string]*http.Cookie

func (c CookieMap) GetValue(key string) string {
	if v, ok := c[key]; ok {
		return v.Value
	}
	return ""
}

type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

type Client struct {
	*http.Client
	cookies CookieMap
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Jar: new(Jar),
		},
		cookies: make(CookieMap),
	}
}

func NewProxyClient(proxyURL string) (*Client, error) {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: &http.Client{
			Jar: new(Jar),
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		},
		cookies: make(CookieMap),
	}, nil
}

func (cli *Client) IsExists(url string) bool {
	resp, err := cli.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return false
	}
	return true
}
