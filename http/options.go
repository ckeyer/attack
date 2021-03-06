package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ckeyer/api/types"
)

// CheckOption check types.HTTPOption.
func CheckOption(opt types.HTTPOption) error {
	_, err := url.Parse(opt.Url)
	if err != nil {
		return err
	}

	switch opt.Method {
	case "GET", "POST", "PUT", "PATCH", "HEAD", "DELETE":
	default:
		return fmt.Errorf("not support http method %s", opt.Method)
	}

	return nil
}

// GetOptionHeader get header from option
func GetOptionHeader(opt types.HTTPOption) http.Header {
	hdr := http.Header{}
	for k, v := range opt.Headers {
		hdr.Add(k, v)
	}
	return hdr
}
