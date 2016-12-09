package juemei

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/ckeyer/attack/httpclient"
)

func LoadDoc(cli *httpclient.Client, Url string) (*goquery.Document, error) {
	resp, err := cli.Get(Url)
	if err != nil {
		log.Errorf("GET failed, %s %v", Url, err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("New goquery Document failed, %v", err)
		return nil, err
	}

	doc.Url, _ = url.Parse(Url)
	return doc, nil
}

func Resolve(cli *httpclient.Client, doc *goquery.Document) (*ResolveResult, error) {
	imgUrls := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imgUrl := s.AttrOr("src", "")
		imgUrl = strings.Replace(imgUrl, "_s.jpg", ".jpg", -1)
		if cli.IsExists(imgUrl) {
			imgUrls = append(imgUrls, imgUrl)
		}
	})

	links := []string{}
	outLinks := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if !ok || cli.IsExists(link) {
			return
		}

		if strings.HasPrefix(link, "/") {
			links = append(links, RootURL+link)
			return
		}

		for _, ignore := range IgnorePrefix {
			if strings.HasPrefix(link, ignore) {
				return
			}
		}

		for _, black := range BlackPrefix {
			if strings.HasPrefix(link, black) {
				links = append(links, link)
				return
			}
		}
		outLinks = append(outLinks, link)
	})

	return &ResolveResult{
		URL:      doc.Url.String(),
		Imgs:     imgUrls,
		Links:    links,
		OutLinks: outLinks,
	}, nil
}
