package greeting

import (
	"time"
)

type FetcherMock struct {
	MockFetchIndexPage       func() (*IndexPage, error)
	MockFetchSecretIndexPage func(date time.Time) (*IndexPage, error)
	MockFetchMenuPage        func(path string) (*MenuPage, error)
	MockFetchCharacterPage   func(path string) (*CharacterPage, error)
}

func (f *FetcherMock) FetchIndexPage() (*IndexPage, error) {
	return f.MockFetchIndexPage()
}

func (f *FetcherMock) FetchSecretIndexPage(date time.Time) (*IndexPage, error) {
	return f.MockFetchSecretIndexPage(date)
}

func (f *FetcherMock) FetchMenuPage(path string) (*MenuPage, error) {
	return f.MockFetchMenuPage(path)
}

func (f *FetcherMock) FetchCharacterPage(path string) (*CharacterPage, error) {
	return f.MockFetchCharacterPage(path)
}
