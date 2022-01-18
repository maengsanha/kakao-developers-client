package local_test

import (
	"testing"

	"github.com/maengsanha/kakao-developers-api/local"
)

func TestCategorySearchWithJSON(t *testing.T) {
	key := ""
	var x float64 = 127.06283102249932
	var y float64 = 37.514322572335935
	radius := 2000
	groupcode := "MT1"

	iter := local.PlaceSearchByCategory(groupcode).
		FormatJSON().
		AuthorizeWith(key).
		WithRadius(x, y, radius).
		Display(15).
		Result(1)

	for res, err := iter.Next(); ; res, err = iter.Next() {
		t.Log(res)
		if err != nil {
			if err != local.ErrEndPage {
				t.Error(err)
			}
			break
		}
	}
}

func TestCategorySearchWithXML(t *testing.T) {
	key := ""
	groupcode := "CS2"
	xmin := 127.05897078335246
	ymin := 37.506051888130386
	xmax := 128.05897078335276
	ymax := 38.506051888130406

	iter := local.PlaceSearchByCategory(groupcode).
		FormatXML().
		AuthorizeWith(key).
		WithRect(xmin, ymin, xmax, ymax).
		Display(15).
		Result(1)

	for res, err := iter.Next(); ; res, err = iter.Next() {
		t.Log(res)
		if err != nil {
			if err != local.ErrEndPage {
				t.Error(err)
			}
			break
		}
	}
}