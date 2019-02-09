package greeting

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWalker(t *testing.T) {
	fetcher := &FetcherMock{}
	walker := NewWalker(fetcher, 0)
	assert.NotNil(t, walker)
}

func TestWalk(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	fetcher := &FetcherMock{}
	fetcher.MockFetchIndexPage = func() (*IndexPage, error) {
		return &IndexPage{
			Date:         time.Date(2018, 10, 20, 0, 0, 0, 0, loc),
			MenuPagePath: "/path-to-menu-page",
		}, nil
	}
	fetcher.MockFetchMenuPage = func(path string) (*MenuPage, error) {
		assert.Equal(t, "/path-to-menu-page", path)
		return &MenuPage{
			Items: []*MenuPageItem{
				{
					CharacterName:     "キティ・ホワイト",
					CharacterPagePath: "/path-to-character-page-kitty",
				},
				{
					CharacterName:     "シナモン",
					CharacterPagePath: "/path-to-character-page-cinnamon",
				},
			},
		}, nil
	}
	fetcher.MockFetchCharacterPage = func(path string) (*CharacterPage, error) {
		return map[string]*CharacterPage{
			"/path-to-character-page-kitty": {
				Name: "キティ・ホワイト",
				Date: time.Date(2018, 10, 20, 0, 0, 0, 0, loc),
				Items: []*CharacterPageItem{
					{
						StartTime: time.Date(2018, 10, 20, 10, 20, 0, 0, loc),
						EndTime:   time.Date(2018, 10, 20, 10, 50, 0, 0, loc),
						Venue:     "ビレッジ(1F)",
					},
					{
						StartTime: time.Date(2018, 10, 20, 15, 0, 0, 0, loc),
						EndTime:   time.Date(2018, 10, 20, 15, 30, 0, 0, loc),
						Venue:     "館外(他)",
					},
				},
			},
			"/path-to-character-page-cinnamon": {
				Name: "シナモン",
				Date: time.Date(2018, 10, 20, 0, 0, 0, 0, loc),
				Items: []*CharacterPageItem{
					{
						StartTime: time.Date(2018, 10, 20, 10, 20, 0, 0, loc),
						EndTime:   time.Date(2018, 10, 20, 10, 50, 0, 0, loc),
						Venue:     "ビレッジ(1F)",
					},
					{
						StartTime: time.Date(2018, 10, 20, 13, 30, 0, 0, loc),
						EndTime:   time.Date(2018, 10, 20, 14, 0, 0, 0, loc),
						Venue:     "4Fふわもこタウンシナモロールわごん付近(4F)",
					},
				},
			},
		}[path], nil
	}

	walker := NewWalker(fetcher, 0)
	rawData, err := walker.Walk()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Date(2018, 10, 20, 0, 0, 0, 0, loc), rawData.IndexPage.Date)
	assert.Equal(t, "/path-to-menu-page", rawData.IndexPage.MenuPagePath)

	assert.Equal(t, 2, len(rawData.MenuPage.Items))
	assert.Equal(t, "キティ・ホワイト", rawData.MenuPage.Items[0].CharacterName)
	assert.Equal(t, "/path-to-character-page-kitty", rawData.MenuPage.Items[0].CharacterPagePath)

	assert.Equal(t, 2, len(rawData.CharacterPages))
	assert.Equal(t, "キティ・ホワイト", rawData.CharacterPages[0].Name)
	assert.Equal(t, time.Date(2018, 10, 20, 0, 0, 0, 0, loc), rawData.CharacterPages[0].Date)
	assert.Equal(t, 2, len(rawData.CharacterPages[0].Items))
	assert.Equal(t, time.Date(2018, 10, 20, 10, 20, 0, 0, loc), rawData.CharacterPages[0].Items[0].StartTime)
	assert.Equal(t, time.Date(2018, 10, 20, 10, 50, 0, 0, loc), rawData.CharacterPages[0].Items[0].EndTime)
	assert.Equal(t, "ビレッジ(1F)", rawData.CharacterPages[0].Items[0].Venue)
}

func TestWalk_unpublished(t *testing.T) {
	fetcher := &FetcherMock{}
	fetcher.MockFetchIndexPage = func() (*IndexPage, error) {
		return nil, errors.New("not found")
	}

	walker := NewWalker(fetcher, 0)
	rawData, err := walker.Walk()

	assert.Nil(t, rawData)
	assert.NotNil(t, err)
}

func TestWalk_secret(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	fetcher := &FetcherMock{}
	fetcher.MockFetchIndexPage = func() (*IndexPage, error) {
		return nil, &SecretError{Date: time.Date(2018, 6, 28, 0, 0, 0, 0, loc)}
	}
	fetcher.MockFetchSecretIndexPage = func(date time.Time) (*IndexPage, error) {
		assert.Equal(t, time.Date(2018, 6, 28, 0, 0, 0, 0, loc), date)
		return &IndexPage{
			Date:         time.Date(2018, 6, 28, 0, 0, 0, 0, loc),
			MenuPagePath: "/path-to-menu-page",
		}, nil
	}
	fetcher.MockFetchMenuPage = func(path string) (*MenuPage, error) {
		assert.Equal(t, "/path-to-menu-page", path)
		return &MenuPage{
			Items: []*MenuPageItem{
				{
					CharacterName:     "キティ・ホワイト",
					CharacterPagePath: "/path-to-character-page-kitty",
				},
			},
		}, nil
	}
	fetcher.MockFetchCharacterPage = func(path string) (*CharacterPage, error) {
		return map[string]*CharacterPage{
			"/path-to-character-page-kitty": {
				Name: "キティ・ホワイト",
				Date: time.Date(2018, 6, 28, 0, 0, 0, 0, loc),
				Items: []*CharacterPageItem{
					{
						StartTime: time.Date(2018, 6, 28, 10, 20, 0, 0, loc),
						EndTime:   time.Date(2018, 6, 28, 10, 50, 0, 0, loc),
						Venue:     "ビレッジ(1F)",
					},
				},
			},
		}[path], nil
	}

	walker := NewWalker(fetcher, 0)
	rawData, err := walker.Walk()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Date(2018, 6, 28, 0, 0, 0, 0, loc), rawData.IndexPage.Date)
	assert.Equal(t, "/path-to-menu-page", rawData.IndexPage.MenuPagePath)

	assert.Equal(t, 1, len(rawData.MenuPage.Items))
	assert.Equal(t, "キティ・ホワイト", rawData.MenuPage.Items[0].CharacterName)
	assert.Equal(t, "/path-to-character-page-kitty", rawData.MenuPage.Items[0].CharacterPagePath)

	assert.Equal(t, 1, len(rawData.CharacterPages))
	assert.Equal(t, "キティ・ホワイト", rawData.CharacterPages[0].Name)
	assert.Equal(t, time.Date(2018, 6, 28, 0, 0, 0, 0, loc), rawData.CharacterPages[0].Date)
	assert.Equal(t, 1, len(rawData.CharacterPages[0].Items))
	assert.Equal(t, time.Date(2018, 6, 28, 10, 20, 0, 0, loc), rawData.CharacterPages[0].Items[0].StartTime)
	assert.Equal(t, time.Date(2018, 6, 28, 10, 50, 0, 0, loc), rawData.CharacterPages[0].Items[0].EndTime)
	assert.Equal(t, "ビレッジ(1F)", rawData.CharacterPages[0].Items[0].Venue)
}
