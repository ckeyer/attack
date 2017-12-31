package name

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	RootURL = "http://www.resgain.net/"
)

func alphaURL(abc byte) string {
	return fmt.Sprintf("%sxmdq_%s.html", RootURL, []byte{abc})
}

func Download(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	doc.Find(".panel-info .panel-body a").Not(".row .fmlink a").Each(func(i int, e *goquery.Selection) {
		ret = append(ret, e.Text())
	})

	return ret, nil
}

func Update() error {
	os.Remove("names.txt")
	f, err := os.OpenFile("names.txt", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewWriter(f)

	for i := byte('a'); i <= 'z'; i++ {
		ss, err := Download(alphaURL(i))
		if err != nil {
			return err
		}
		r := strings.NewReader(strings.Join(ss, "\n"))

		_, err = io.Copy(buf, r)
		if err != nil {
			return err
		}
	}
	return nil
}
