package internal_test

import (
	"github.com/provokateurin/rwth-calendar/internal"
	"testing"
)

var h05 = &internal.RWTHLocation{
	BuildingID: stringP("1385"),
	RoomID:     stringP("105"),
	Name:       "H05",
}

func TestNewRWTHLocation(t *testing.T) {
	t.Parallel()

	cases := map[string]*internal.RWTHLocation{
		"[1385|105] H05": h05,
		"H05 (1385|105)": h05,
	}

	for location, expected := range cases {
		actual := internal.NewRWTHLocation(location)
		if !actual.Equal(expected) {
			t.Errorf("Parsing location from %q was not successful", location)
		}
	}
}

func stringP(s string) *string {
	return &s
}
