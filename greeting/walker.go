package greeting

import (
	"time"
)

type Walker interface {
	Walk() (*RawData, error)
}

type walkerImpl struct {
	fetcher  Fetcher
	duration time.Duration
}

func NewWalker(fetcher Fetcher, duration time.Duration) Walker {
	return &walkerImpl{
		fetcher:  fetcher,
		duration: duration,
	}
}

func (w *walkerImpl) sleep() {
	if w.duration > 0 {
		time.Sleep(w.duration)
	}
}

func (w *walkerImpl) Walk() (*RawData, error) {
	indexPage, err := w.fetcher.FetchIndexPage()
	if err != nil {
		secretErr, ok := err.(*SecretError)
		if !ok {
			return nil, err
		}

		w.sleep()

		indexPage, err = w.fetcher.FetchSecretIndexPage(secretErr.Date)
		if err != nil {
			return nil, err
		}
	}

	w.sleep()

	menuPage, err := w.fetcher.FetchMenuPage(indexPage.MenuPagePath)
	if err != nil {
		return nil, err
	}

	characterPages := make([]*CharacterPage, 0, len(menuPage.Items))
	for _, item := range menuPage.Items {
		w.sleep()

		characterPage, err := w.fetcher.FetchCharacterPage(item.CharacterPagePath)
		if err != nil {
			return nil, err
		}
		characterPages = append(characterPages, characterPage)
	}
	return &RawData{
		IndexPage:      indexPage,
		MenuPage:       menuPage,
		CharacterPages: characterPages,
	}, nil
}
