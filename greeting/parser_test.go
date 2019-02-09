package greeting

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseIndexPage(t *testing.T) {
	f, err := os.Open("testdata/www.puroland.co.jp/chara_gre/mobile/index.asp")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	parser := NewParser()

	page, err := parser.ParseIndexPage(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "chara_sentaku.asp?TCHK=2016156&lang=", page.MenuPagePath)
	assert.Equal(t, time.Date(2016, time.June, 15, 0, 0, 0, 0, loc), page.Date)
}

func TestParseMenuPage(t *testing.T) {
	f, err := os.Open("testdata/www.puroland.co.jp/chara_gre/mobile/chara_sentaku.asp")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	parser := NewParser()

	page, err := parser.ParseMenuPage(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 12, len(page.Items))
	assert.Equal(t, "キティ・ホワイト", page.Items[0].CharacterName)
	assert.Equal(t, "chara_sche.asp?TCHK=2016156&C_KEY=1", page.Items[0].CharacterPagePath)
}

func TestParseCharacterPage(t *testing.T) {
	f, err := os.Open("testdata/www.puroland.co.jp/chara_gre/mobile/chara_sche.asp")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	parser := NewParser()

	page, err := parser.ParseCharacterPage(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Date(2016, time.June, 15, 0, 0, 0, 0, loc), page.Date)
	assert.Equal(t, "キティ・ホワイト", page.Name)
	assert.Equal(t, 2, len(page.Items))
	assert.Equal(t, time.Date(2016, time.June, 15, 11, 0, 0, 0, loc), page.Items[0].StartTime)
	assert.Equal(t, time.Date(2016, time.June, 15, 11, 30, 0, 0, loc), page.Items[0].EndTime)
	assert.Equal(t, "エンターテイメントホール入口（１Ｆ）", page.Items[0].Venue)
}

func TestParseCharacterListPage(t *testing.T) {
	f, err := os.Open("testdata/www.puroland.jp/greeting/schedule/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	parser := NewParser()

	page, err := parser.ParseCharacterListPage(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Date(2017, time.August, 4, 0, 0, 0, 0, loc), page.Date)
	assert.Equal(t, 21, len(page.Items))
	assert.Equal(t, "シナモン", page.Items[0].Name)
}
