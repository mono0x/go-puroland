package greeting

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetcher(t *testing.T) {
	parser := NewParser()
	fetcher := NewFetcher(http.DefaultClient, parser, BaseURL)
	assert.NotNil(t, fetcher)
}

func TestFetchIndexPage(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.puroland.co.jp/chara_gre/mobile/index.asp")
	})
	server := httptest.NewServer(mux)
	client := server.Client()

	parser := NewParser()
	fetcher := NewFetcher(client, parser, server.URL+"/")

	page, err := fetcher.FetchIndexPage()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "chara_sentaku.asp?TCHK=2016156&lang=", page.MenuPagePath)
	assert.Equal(t, time.Date(2016, time.June, 15, 0, 0, 0, 0, loc), page.Date)
}

func TestFetchMenuPage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/chara_sentaku.asp", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "2016156", query.Get("TCHK"))
		http.ServeFile(w, r, "testdata/www.puroland.co.jp/chara_gre/mobile/chara_sentaku.asp")
	})
	server := httptest.NewServer(mux)
	client := server.Client()

	parser := NewParser()
	fetcher := NewFetcher(client, parser, server.URL+"/")

	page, err := fetcher.FetchMenuPage("chara_sentaku.asp?TCHK=2016156&lang=")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 12, len(page.Items))
	assert.Equal(t, "キティ・ホワイト", page.Items[0].CharacterName)
	assert.Equal(t, "chara_sche.asp?TCHK=2016156&C_KEY=1", page.Items[0].CharacterPagePath)
}

func TestFetchCharacterPage(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/chara_sche.asp", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "2016156", query.Get("TCHK"))
		assert.Equal(t, "1", query.Get("C_KEY"))
		http.ServeFile(w, r, "testdata/www.puroland.co.jp/chara_gre/mobile/chara_sche.asp")
	})
	server := httptest.NewServer(mux)
	client := server.Client()

	parser := NewParser()
	fetcher := NewFetcher(client, parser, server.URL+"/")

	page, err := fetcher.FetchCharacterPage("chara_sche.asp?TCHK=2016156&C_KEY=1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, time.Date(2016, time.June, 15, 0, 0, 0, 0, loc), page.Date)
	assert.Equal(t, "キティ・ホワイト", page.Name)
	assert.Equal(t, 2, len(page.Items))
	assert.Equal(t, time.Date(2016, time.June, 15, 11, 0, 0, 0, loc), page.Items[0].StartTime)
	assert.Equal(t, time.Date(2016, time.June, 15, 11, 30, 0, 0, loc), page.Items[0].EndTime)
}
