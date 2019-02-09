package greeting

import (
	"time"
)

type IndexPage struct {
	Date         time.Time
	MenuPagePath string
}

type MenuPage struct {
	Items []*MenuPageItem
}

type MenuPageItem struct {
	CharacterName     string
	CharacterPagePath string
}

type CharacterPage struct {
	Name  string
	Date  time.Time
	Items []*CharacterPageItem
}

type CharacterPageItem struct {
	StartTime time.Time
	EndTime   time.Time
	Venue     string
}

type CharacterListPage struct {
	Date  time.Time
	Items []*CharacterListPageItem
}

type CharacterListPageItem struct {
	Name string
}
