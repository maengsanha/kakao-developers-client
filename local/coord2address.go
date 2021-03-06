// Copyright 2022 Sanha Maeng, Soyang Baek, Jinmyeong Kim
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package local

import (
	"encoding/xml"
	"errors"
	"fmt"
	"internal/common"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

// TotalAddress represents a document of CoordToAddressResult.
type TotalAddress struct {
	Address struct {
		AddressName      string `json:"address_name" xml:"address_name"`
		Region1depthName string `json:"region_1depth_name" xml:"region_1depth_name"`
		Region2depthName string `json:"region_2depth_name" xml:"region_2depth_name"`
		Region3depthName string `json:"region_3depth_name" xml:"region_3depth_name"`
		MountainYN       string `json:"mountain_yn" xml:"mountain_yn"`
		MainAddressNo    string `json:"main_address_no" xml:"main_address_no"`
		SubAddressNo     string `json:"sub_address_no" xml:"sub_address_no"`
		ZipCode          string `json:"zip_code" xml:"zip_code"`
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
	} `json:"road_address" xml:"road_address"`
}

// CoordToAddressResult represents a CoordToAddress result.
type CoordToAddressResult struct {
	XMLName   xml.Name       `json:"-" xml:"result"`
	Meta      common.Meta    `json:"meta" xml:"meta"`
	Documents []TotalAddress `json:"documents" xml:"documents"`
}

// String implements fmt.Stringer.
func (cr CoordToAddressResult) String() string { return common.String(cr) }

// SaveAs saves cr to @filename.
//
// The file extension could be either .json or .xml.
func (cr CoordToAddressResult) SaveAs(filename string) error {
	return common.SaveAsJSONorXML(cr, filename)
}

// CoordToAddressInitializer is a lazy coord to address converter.
type CoordToAddressInitializer struct {
	X          string
	Y          string
	Format     string
	AuthKey    string
	InputCoord string
}

// CoordToAddress converts the @x and @y coordinates of location in the selected coordinate system
// to land-lot number address(with post number) and road name address.
//
// Details can be referred to
// https://developers.kakao.com/docs/latest/ko/local/dev-guide#coord-to-address.
func CoordToAddress(x, y string) *CoordToAddressInitializer {
	return &CoordToAddressInitializer{
		X:          x,
		Y:          y,
		Format:     "json",
		AuthKey:    common.KeyPrefix,
		InputCoord: "WGS84",
	}
}

// FormatAs sets the request format to @format (json or xml).
func (ci *CoordToAddressInitializer) FormatAs(format string) *CoordToAddressInitializer {
	switch format {
	case "json", "xml":
		ci.Format = format
	default:
		panic(common.ErrUnsupportedFormat)
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return ci
}

// AuthorizeWith sets the authorization key to @key.
func (ci *CoordToAddressInitializer) AuthorizeWith(key string) *CoordToAddressInitializer {
	ci.AuthKey = common.FormatKey(key)
	return ci
}

// Input sets the coordinate system of request.
//
// There are following coordinate system exist:
//
// WGS84
//
// WCONGNAMUL
//
// CONGNAMUL
//
// WTM
//
// TM
func (ci *CoordToAddressInitializer) Input(coord string) *CoordToAddressInitializer {
	switch coord {
	case "WGS84", "WCONAMUL", "CONGNAMUL", "WTM", "TM":
		ci.InputCoord = coord
	default:
		panic(errors.New(
			`input coordinate system must be one of following options:
			WGS84, WCONGNAMUL, CONGNAMUL, WTM, TM`))
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return ci
}

// Collect returns the land-lot number address(with post number) and road name address.
func (ci *CoordToAddressInitializer) Collect() (res CoordToAddressResult, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%sgeo/coord2address.%s?x=%s&y=%s&input_coord=%s",
			prefix, ci.Format, ci.X, ci.Y, ci.InputCoord), nil)

	if err != nil {
		return
	}

	req.Close = true

	req.Header.Set(common.Authorization, ci.AuthKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if ci.Format == "json" {
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	} else if ci.Format == "xml" {
		if err = xml.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	}

	return
}
