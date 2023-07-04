package internal

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"sync"
)

type RWTHBuildingDetails struct {
	ID     string
	Name   string
	Street string
	Place  string
}

var (
	buildings *[]RWTHBuildingDetails
	lock      = &sync.Mutex{}
)

func NewRWTHBuildingDetails(buildingID string) (*RWTHBuildingDetails, error) {
	err := loadBuildingsTable()
	if err != nil {
		return nil, err
	}

	for _, building := range *buildings {
		if building.ID == buildingID {
			return &building, nil
		}
	}

	return nil, nil
}

func loadBuildingsTable() error {
	if buildings != nil {
		return nil
	}

	lock.Lock()

	buildingsResponse, err := http.Get("https://www.rwth-aachen.de/cms/root/Die-RWTH/Kontakt-Anreise/RWTH-Navigator/~cxcq/Maps-Gebaeude/?showall=1")
	if err != nil {
		return err
	}
	defer buildingsResponse.Body.Close()

	buildingsResponseBody, err := io.ReadAll(buildingsResponse.Body)
	if err != nil {
		return err
	}

	buildingsDoc, err := goquery.NewDocumentFromReader(strings.NewReader(iso8859toUtf8(buildingsResponseBody)))
	if err != nil {
		return err
	}

	buildings = &[]RWTHBuildingDetails{}

	buildingsDoc.Find(".mod").Children().Slice(1, goquery.ToEnd).Each(func(i int, selection *goquery.Selection) {
		elements := selection.Children().Slice(1, goquery.ToEnd)
		*buildings = append(*buildings, RWTHBuildingDetails{
			ID:     elements.Eq(0).Text(),
			Name:   elements.Eq(1).Text(),
			Street: elements.Eq(2).Text(),
			Place:  elements.Eq(3).Text(),
		})
	})

	lock.Unlock()

	return nil
}

func iso8859toUtf8(bytes []byte) string {
	buf := make([]rune, len(bytes))

	for i, b := range bytes {
		buf[i] = rune(b)
	}

	return string(buf)
}
