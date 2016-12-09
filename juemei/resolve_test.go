package juemei

import (
	"testing"

	"github.com/ckeyer/attack/httpclient"
)

func TestGetPage(t *testing.T) {
	url := "http://www.juemei.com/mm/201612/6386.html"
	cli := httpclient.NewClient()
	doc, err := LoadDoc(cli, url)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := Resolve(cli, doc)
	if err != nil {
		t.Fatal(err)
	}

	if ret.URL != url {
		t.Fatal("url not equie")
	}
}
