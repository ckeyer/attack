package blockai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ckeyer/attack/httpclient"
)

const (
	URLAuth        = "https://blockai.com/api/v1/auth/"
	URLInitProfile = "https://blockai.com/api/v1/profile"
	URLUpload      = "https://files.blockai.com/v1/upload"
	URLReg         = "https://blockai.com/api/v1/registrations"
)

type Img struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	IsPrivate bool   `json:"is_private"`
}

type Client struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	token    string

	cli *httpclient.Client
}

func NewClient(name, email, passwd string) (*Client, error) {
	cli := &Client{
		Name:     name,
		Email:    email,
		Password: passwd,
		cli:      httpclient.NewClient(),
	}
	if err := cli.signupAndlogin(); err != nil {
		return nil, err
	}

	profile := map[string]map[string]interface{}{
		"ui": map[string]interface{}{
			"introMessageShowed": true,
		},
	}
	_, err := cli.PostJSON(URLInitProfile, profile)
	return cli, err
}

func NewRandomClient() (*Client, error) {
	return nil, nil
}

func (i *Img) FmtName() {
	ps := strings.Split(i.URL, "/")
	f := ps[len(ps)-1]
	f = strings.Replace(f, "-", " ", -1)
	f = strings.Replace(f, "_", " ", -1)
	index := strings.LastIndex(f, ".")
	f = string(f[:index])
	i.Name = strings.Title(strings.ToLower(f))
}

func (c *Client) signupAndlogin() error {
	for _, action := range []string{"signup", "login"} {
		bs, err := json.Marshal(c)
		if err != nil {
			return err
		}

		resp, err := c.cli.Post(URLAuth+action, "application/json", bytes.NewReader(bs))
		if err != nil {
			return err
		}

		ret := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&ret)
		if err != nil {
			return fmt.Errorf("%s failed, %s, body: %+v", action, resp.Status, ret)
		}
		if tok, ok := ret["token"]; ok {
			c.token = tok.(string)
		}

		log.Infof("%s a account %+v, return: %+v", action, c.Email, ret)
	}
	return nil
}

func (c *Client) Upload(img *Img) error {
	if err := c.uploadImg(img); err != nil {
		return err
	}
	if err := c.registrations(img); err != nil {
		return err
	}
	return nil
}

func (c *Client) uploadImg(img *Img) error {
	resp, err := c.PostJSON(URLUpload, img)
	if err != nil {
		return err
	}

	ret := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}

	if url, ok := ret["url"]; ok {
		img.URL = url.(string)
	}

	return nil
}

func (c *Client) registrations(img *Img) error {
	resp, err := c.PostJSON(URLReg, img)
	if err != nil {
		return err
	}

	ret := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Errorf("%+v", ret)
		return err
	}

	log.Infof("registrations, %+v", ret)
	return nil
}

func (c *Client) PostJSON(url string, data interface{}) (*http.Response, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-identity-token", c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.cli.Do(req)
}
