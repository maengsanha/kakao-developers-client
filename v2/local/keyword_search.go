// Package local provides the features of the Local API.
package local

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultCategoryGroupCode = ""
	DefaultX                 = ""
	DefaultY                 = ""
	DefaultRadius            = 0
	MaxRadius                = 20000
	MinRadius                = 0
	DefaultRect              = ""
	DefaultSort              = "accuracy"
	Distance                 = "distance"
	Accuracy                 = "accuracy"
)

//category_group_code_list := []string{"MT1", "CS2", "PS3", "SC4", "AC5", "PK6", "OL7", "SW8" ,"BK9", "CT1", "AG2", "PO3", "AT4", "AD5", "FD6", "CE7", "HP8", "PM9", ""}

type SameName struct {
	Region         []string `json:"region" xml:"region"`
	Keyword        string   `json:"keyword" xml:"keyword"`
	SelectedRegion string   `json:"selected_region" xml:"selected_region"`
}

type Place struct {
	Id                string `json:"id" xml:"id"`
	PlaceName         string `json:"place_name" xml:"place_name"`
	CategoryName      string `json:"category_name" xml:"category_name"`
	CategoryGroupCode string `json:"category_group_code" xml:"category_group_code"`
	CategoryGroupName string `json:"category_group_name" xml:"category_group_name"`
	Phone             string `json:"phone" xml:"phone"`
	AddressName       string `json:"address_name" xml:"address_name"`
	RoadAddressName   string `json:"road_address_name" xml:"road_address_name"`
	X                 string `json:"x" xml:"x"`
	Y                 string `json:"y" xml:"y"`
	PlaceURL          string `json:"place_url" xml:"place_url"`
	Distance          string `json:"distance" xml:"distance"`
}

type KeyWordSearchResult struct {
	XMLName xml.Name `xml:"result"`
	Meta    struct {
		TotalCount    int  `json:"total_count" xml:"total_count"`
		PageableCount int  `json:"pageable_count" xml:"pageable_count"`
		IsEnd         bool `json:"is_end" xml:"is_end"`
		SameName      `json:"same_name" xml:"same_name"`
	} `json:"meta" xml:"meta"`
	Documents []Place `json:"documents" xml:"documents"`
}

type KeyWordSearchInitializer struct {
	Query             string
	CategoryGroupCode string
	Format            string
	AuthKey           string
	X                 string
	Y                 string
	Radius            int
	Rect              string
	Page              int
	Size              int
	Sort              string
}

func KeyWordSearch(query string) *KeyWordSearchInitializer {
	return &KeyWordSearchInitializer{
		Query:             url.QueryEscape(strings.TrimSpace(query)),
		CategoryGroupCode: DefaultCategoryGroupCode,
		Format:            JSON,
		AuthKey:           keyPrefix,
		X:                 DefaultX,
		Y:                 DefaultY,
		Radius:            DefaultRadius,
		Rect:              DefaultRect,
		Page:              1,
		Size:              15,
		Sort:              DefaultSort,
	}
}

func (k *KeyWordSearchInitializer) As(format string) *KeyWordSearchInitializer {
	if format == JSON || format == XML {
		k.Format = format
	}
	return k
}

func (k *KeyWordSearchInitializer) AuthorizeWith(key string) *KeyWordSearchInitializer {
	k.AuthKey = keyPrefix + strings.TrimSpace(key)
	return k
}

func (k *KeyWordSearchInitializer) SetCategoryGroupCode(code string) *KeyWordSearchInitializer {
	k.CategoryGroupCode = code
	return k
}

func (k *KeyWordSearchInitializer) SetX(x string) *KeyWordSearchInitializer {
	k.X = x
	return k
}

func (k *KeyWordSearchInitializer) SetY(y string) *KeyWordSearchInitializer {
	k.Y = y
	return k
}

func (k *KeyWordSearchInitializer) SetRadius(radius int) *KeyWordSearchInitializer {
	if MinRadius <= radius && radius <= MaxRadius {
		k.Radius = radius
	}
	return k
}

func (k *KeyWordSearchInitializer) SetRect(rect string) *KeyWordSearchInitializer {
	k.Rect = rect

	return k
}

func (k *KeyWordSearchInitializer) Result(page int) *KeyWordSearchInitializer {
	if 1 <= page && page <= 45 {
		k.Page = page
	}
	return k
}

func (k *KeyWordSearchInitializer) Display(size int) *KeyWordSearchInitializer {
	if 1 <= size && size <= 45 {
		k.Size = size
	}
	return k
}

func (k *KeyWordSearchInitializer) SortType(sort string) *KeyWordSearchInitializer {
	if sort == Accuracy || sort == Distance {
		k.Sort = sort
	}
	return k
}

// Next returns the search result and proceeds the iterator to the next page.
func (k *KeyWordSearchInitializer) Collect() (res KeyWordSearchResult, err error) {
	// at first, send request to the API server
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.%s?query=%s&category_group_code=%s&x=%s&y=%s&radius=%d&rect=%s&page=%d&size=%d&sort=%s", k.Format, k.Query, k.CategoryGroupCode, k.X, k.Y, k.Radius, k.Rect, k.Page, k.Size, k.Sort), nil)
	if err != nil {
		return
	}
	// don't forget to close the request for concurrent request
	req.Close = true

	// set authorization header
	req.Header.Set(authorization, k.AuthKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// don't forget to close the response body
	defer resp.Body.Close()

	if k.Format == JSON {
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	} else if k.Format == XML {
		if err = xml.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	}

	return
}
