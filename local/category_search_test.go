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

package local_test

import (
	"internal/common"
	"testing"

	"github.com/maengsanha/kakao-developers-client/local"
)

func TestCategorySearchWithJSON(t *testing.T) {
	var x float64 = 127.06283102249932
	var y float64 = 37.514322572335935
	radius := 2000
	groupcode := "MT1"

	it := local.PlaceSearchByCategory(groupcode).
		FormatAs("json").
		AuthorizeWith(common.REST_API_KEY).
		WithRadius(x, y, radius).
		Display(15).
		Result(1)

	for {
		item, err := it.Next()
		if err == local.Done {
			break
		}
		if err != nil {
			t.Error(err)
		}
		t.Log(item)
	}
}

func TestCategorySearchWithSaveAsJSON(t *testing.T) {
	var x float64 = 127.06283102249932
	var y float64 = 37.514322572335935
	radius := 2000
	groupcode := "MT1"

	it := local.PlaceSearchByCategory(groupcode).
		FormatAs("json").
		AuthorizeWith(common.REST_API_KEY).
		WithRadius(x, y, radius).
		Display(15).
		Result(1)

	items := local.PlaceSearchResults{}

	for {
		item, err := it.Next()
		if err == local.Done {
			break
		}
		if err != nil {
			t.Error(err)
		}
		items = append(items, item)
	}

	if err := items.SaveAs("category_search_test.json"); err != nil {
		t.Error(err)
	}
}

func TestCategorySearchWithXML(t *testing.T) {
	groupcode := "CS2"
	xmin := 127.05897078335246
	ymin := 37.506051888130386
	xmax := 128.05897078335276
	ymax := 38.506051888130406

	it := local.PlaceSearchByCategory(groupcode).
		FormatAs("xml").
		AuthorizeWith(common.REST_API_KEY).
		WithRect(xmin, ymin, xmax, ymax).
		Display(15).
		Result(1)

	for {
		item, err := it.Next()
		if err == local.Done {
			break
		}
		if err != nil {
			t.Error(err)
		}
		t.Log(item)
	}
}

func TestCategorySearchWithSaveAsXML(t *testing.T) {
	var x float64 = 127.06283102249932
	var y float64 = 37.514322572335935
	radius := 2000
	groupcode := "MT1"

	it := local.PlaceSearchByCategory(groupcode).
		FormatAs("xml").
		AuthorizeWith(common.REST_API_KEY).
		WithRadius(x, y, radius).
		Display(15).
		Result(1)

	items := local.PlaceSearchResults{}

	for {
		item, err := it.Next()
		if err == local.Done {
			break
		}
		if err != nil {
			t.Error(err)
		}
		items = append(items, item)
	}

	if err := items.SaveAs("category_search_test.xml"); err != nil {
		t.Error(err)
	}
}

func TestCategorySearchCollectAll(t *testing.T) {
	var x float64 = 127.06283102249932
	var y float64 = 37.514322572335935
	radius := 2000
	groupcode := "MT1"

	items := local.PlaceSearchByCategory(groupcode).
		FormatAs("xml").
		AuthorizeWith(common.REST_API_KEY).
		WithRadius(x, y, radius).
		Display(15).
		Result(1).
		CollectAll()

	for _, item := range items {
		t.Log(item)
	}
}
