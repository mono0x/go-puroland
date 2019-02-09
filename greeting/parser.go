package greeting

import (
	"errors"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/PuerkitoBio/goquery"
)

var (
	dateRe         = regexp.MustCompile(`(\d+)年(\d+)月(\d+)日(?:\([日月火水木金土]\))?`)
	timeAndVenueRe = regexp.MustCompile(`\A\s*([０-９]{1,2})：([０-９]{1,2})－([０-９]{1,2})：([０-９]{1,2})(.+)\s*\z`)
)

type Parser interface {
	ParseIndexPage(r io.Reader) (*IndexPage, error)
	ParseMenuPage(r io.Reader) (*MenuPage, error)
	ParseCharacterPage(r io.Reader) (*CharacterPage, error)
	ParseCharacterListPage(r io.Reader) (*CharacterListPage, error)
}

type SecretError struct {
	Date time.Time
}

func (err *SecretError) Error() string {
	return "" // unused
}

type parserImpl struct {
}

func NewParser() Parser {
	return &parserImpl{}
}

func (p *parserImpl) ParseIndexPage(r io.Reader) (*IndexPage, error) {
	decodedReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}

	date, err := p.parseDate(doc.Find(`p[align="center"] font[size="-1"]`).First().Text())
	if err != nil {
		return nil, err
	}

	secret := false
	doc.Find("p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "本日のｷｬﾗｸﾀｰ情報は公開されておりません。P") {
			secret = true
			return false // break
		}
		return true
	})
	if secret {
		return nil, &SecretError{Date: date}
	}

	form := doc.Find("form").First()
	values := url.Values{}
	form.Find("input").Each(func(_ int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if !exists {
			return
		}
		value, exists := s.Attr("value")
		if !exists {
			return
		}
		values.Add(name, value)
	})

	return &IndexPage{
		MenuPagePath: form.AttrOr("action", "") + "?" + values.Encode(),
		Date:         date,
	}, nil
}

func (p *parserImpl) ParseMenuPage(r io.Reader) (*MenuPage, error) {
	decodedReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}

	links := doc.Find(`a[href^="chara_sche.asp?"]`)

	items := make([]*MenuPageItem, 0, links.Size())
	links.Each(func(_ int, s *goquery.Selection) {
		items = append(items, &MenuPageItem{
			CharacterName:     s.Text(),
			CharacterPagePath: s.AttrOr("href", ""),
		})
	})

	return &MenuPage{
		Items: items,
	}, nil
}

func (p *parserImpl) parseDate(s string) (time.Time, error) {
	submatches := dateRe.FindStringSubmatch(s)
	if len(submatches) == 0 {
		return time.Time{}, errors.New("date not found")
	}

	year, err := strconv.Atoi(submatches[1])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(submatches[2])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(submatches[3])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc), nil
}

func (p *parserImpl) ParseCharacterPage(r io.Reader) (*CharacterPage, error) {
	decodedReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}

	date, err := p.parseDate(doc.Find(`p[align="center"] font[size="-1"]`).First().Text())
	if err != nil {
		return nil, err
	}

	name := doc.Find(`p[align="center"] font[size="-1"]`).Eq(1).Text()

	fonts := doc.Find(`p[align="left"] font[size="-1"]`)

	items := make([]*CharacterPageItem, 0, fonts.Size())
	fonts.Each(func(_ int, s *goquery.Selection) {
		submatches := timeAndVenueRe.FindStringSubmatch(s.Text())
		if len(submatches) != 6 {
			return
		}

		startHour, err := strconv.Atoi(norm.NFKC.String(submatches[1]))
		if err != nil {
			return
		}
		startMinute, err := strconv.Atoi(norm.NFKC.String(submatches[2]))
		if err != nil {
			return
		}
		startTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMinute, 0, 0, date.Location())

		endHour, err := strconv.Atoi(norm.NFKC.String(submatches[3]))
		if err != nil {
			return
		}
		endMinute, err := strconv.Atoi(norm.NFKC.String(submatches[4]))
		if err != nil {
			return
		}
		endTime := time.Date(date.Year(), date.Month(), date.Day(), endHour, endMinute, 0, 0, date.Location())

		venue := submatches[5]

		items = append(items, &CharacterPageItem{
			Venue:     venue,
			StartTime: startTime,
			EndTime:   endTime,
		})
	})

	return &CharacterPage{
		Date:  date,
		Name:  name,
		Items: items,
	}, nil
}

func (p *parserImpl) ParseCharacterListPage(r io.Reader) (*CharacterListPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	href, ok := doc.Find(".yesterday a").Attr("href")
	if !ok {
		return nil, errors.New("link not found")
	}

	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	yesterday, err := time.ParseInLocation("20060102", u.Query().Get("date"), loc)
	if err != nil {
		return nil, err
	}

	date := yesterday.AddDate(0, 0, 1)

	charaNames := doc.Find("p.charaName")
	items := make([]*CharacterListPageItem, 0, charaNames.Size())
	charaNames.Each(func(_ int, s *goquery.Selection) {
		items = append(items, &CharacterListPageItem{
			Name: s.Text(),
		})
	})

	return &CharacterListPage{
		Date:  date,
		Items: items,
	}, nil
}
