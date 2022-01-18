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

// ComplexAddress represents a document of address search.
type ComplexAddress struct {
	AddressName string `json:"address_name" xml:"address_name"`
	AddressType string `json:"address_type" xml:"address_type"`
	X           string `json:"x" xml:"x"`
	Y           string `json:"y" xml:"y"`
	Address     struct {
		AddressName       string `json:"address_name" xml:"address_name"`
		Region1depthName  string `json:"region_1depth_name" xml:"region_1depth_name"`
		Region2depthName  string `json:"region_2depth_name" xml:"region_2depth_name"`
		Region3depthName  string `json:"region_3depth_name" xml:"region_3depth_name"`
		Region3depthHName string `json:"region_3depth_h_name" xml:"region_3depth_h_name"`
		HCode             string `json:"h_code" xml:"h_code"`
		BCode             string `json:"b_code" xml:"b_code"`
		MountainYN        string `json:"mountain_yn" xml:"mountain_yn"`
		MainAddressNo     string `json:"main_address_no" xml:"main_address_no"`
		SubAddressNo      string `json:"sub_address_no" xml:"sub_address_no"`
		ZipCode           string `json:"zip_code" xml:"zip_code"`
		X                 string `json:"x" xml:"x"`
		Y                 string `json:"y" xml:"y"`
	} `json:"address" xml:"address"`
	RoadAddress struct {
		AddressName      string `json:"address_name" xml:"address_name"`
		Region1depthName string `json:"region_1depth_name" xml:"region_1depth_name"`
		Region2depthName string `json:"region_2depth_name" xml:"region_2depth_name"`
		Region3depthName string `json:"region_3depth_name" xml:"region_3depth_name"`
		RoadName         string `json:"road_name" xml:"road_name"`
		UndergroundYN    string `json:"underground_yn" xml:"underground_yn"`
		MainBuildingNo   string `json:"main_building_no" xml:"main_building_no"`
		SubBuildingNo    string `json:"sub_building_no" xml:"sub_building_no"`
		BuildingName     string `json:"building_name" xml:"building_name"`
		ZoneNo           string `json:"zone_no" xml:"zone_no"`
		X                string `json:"x" xml:"x"`
		Y                string `json:"y" xml:"y"`
	} `json:"road_address" xml:"road_address"`
}

// AddressSearchResult represents an address search result.
type AddressSearchResult struct {
	XMLName xml.Name `xml:"result"`
	Meta    struct {
		TotalCount    int  `json:"total_count" xml:"total_count"`
		PageableCount int  `json:"pageable_count" xml:"pageable_count"`
		IsEnd         bool `json:"is_end" xml:"is_end"`
	} `json:"meta" xml:"meta"`
	Documents []ComplexAddress `json:"documents" xml:"documents"`
}

// AddressSearchIterator is a lazy address search iterator.
type AddressSearchIterator struct {
	Query       string
	Format      string
	AuthKey     string
	AnalyzeType string
	Page        int
	Size        int
}

// AddressSearch provides the coordinates of the requested address with @query.
//
// See https://developers.kakao.com/docs/latest/ko/local/dev-guide#address-coord for more details.
func AddressSearch(query string) *AddressSearchIterator {
	return &AddressSearchIterator{
		Query:       url.QueryEscape(strings.TrimSpace(query)),
		Format:      "json",
		AuthKey:     "KakaoAK ",
		AnalyzeType: "similar",
		Page:        1,
		Size:        10,
	}
}

func (a *AddressSearchIterator) FormatJSON() *AddressSearchIterator {
	a.Format = "json"
	return a
}

func (a *AddressSearchIterator) FormatXML() *AddressSearchIterator {
	a.Format = "xml"
	return a
}

// AuthorizeWith sets the authorization key to @key.
func (a *AddressSearchIterator) AuthorizeWith(key string) *AddressSearchIterator {
	a.AuthKey = "KakaoAK " + strings.TrimSpace(key)
	return a
}

// Analyze sets the analyze type to @typ.
//
// There are a few supported analyze types:
//
// similar
//
// exact
func (a *AddressSearchIterator) Analyze(typ string) *AddressSearchIterator {
	switch typ {
	case "similar", "exact":
		a.AnalyzeType = typ
	}
	return a
}

func (a *AddressSearchIterator) Result(page int) *AddressSearchIterator {
	if 1 <= page && page <= 45 {
		a.Page = page
	}
	return a
}

func (a *AddressSearchIterator) Display(size int) *AddressSearchIterator {
	if 1 <= size && size <= 30 {
		a.Size = size
	}
	return a
}

// Next returns the address search result and proceeds the iterator to the next page.
func (a *AddressSearchIterator) Next() (res AddressSearchResult, err error) {
	// at first, send request to the API server
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://dapi.kakao.com/v2/local/search/address.%s?query=%s&analyze_type=%s&page=%d&size=%d",
			a.Format, a.Query, a.AnalyzeType, a.Page, a.Size), nil)

	if err != nil {
		return
	}
	// don't forget to close the request for concurrent request
	req.Close = true

	// set authorization header
	req.Header.Set("Authorization", a.AuthKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// don't forget to close the response body
	defer resp.Body.Close()

	if a.Format == "json" {
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	} else if a.Format == "xml" {
		if err = xml.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	}

	// if it was the last result, return error
	// or increase the page number
	if res.Meta.IsEnd {
		return res, ErrEndPage
	}

	a.Page++

	return
}