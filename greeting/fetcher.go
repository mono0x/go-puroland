package greeting

import (
	"net/http"
	"time"
)

const (
	BaseURL = "http://www.puroland.co.jp/chara_gre/mobile/"
)

type Fetcher interface {
	FetchIndexPage() (*IndexPage, error)
	FetchSecretIndexPage(date time.Time) (*IndexPage, error)
	FetchMenuPage(path string) (*MenuPage, error)
	FetchCharacterPage(path string) (*CharacterPage, error)
}

type fetcherImpl struct {
	client  *http.Client
	parser  Parser
	baseURL string
}

func NewFetcher(client *http.Client, parser Parser, baseURL string) Fetcher {
	return &fetcherImpl{
		client:  client,
		parser:  parser,
		baseURL: baseURL,
	}
}

func (f *fetcherImpl) FetchIndexPage() (*IndexPage, error) {
	res, err := f.client.Get(f.baseURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return f.parser.ParseIndexPage(res.Body)
}

func (f *fetcherImpl) getSecretIndexPageURL(date time.Time) string {
	return f.baseURL + "?para=" + date.Format("20060102")
}

func (f *fetcherImpl) FetchSecretIndexPage(date time.Time) (*IndexPage, error) {
	res, err := f.client.Get(f.getSecretIndexPageURL(date))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return f.parser.ParseIndexPage(res.Body)
}

func (f *fetcherImpl) FetchMenuPage(path string) (*MenuPage, error) {
	res, err := f.client.Get(f.baseURL + path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return f.parser.ParseMenuPage(res.Body)
}

func (f *fetcherImpl) FetchCharacterPage(path string) (*CharacterPage, error) {
	res, err := f.client.Get(f.baseURL + path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return f.parser.ParseCharacterPage(res.Body)
}
