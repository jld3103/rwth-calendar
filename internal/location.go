package internal

import (
	"regexp"
)

var moodleLocationRegexp = regexp.MustCompile("^\\[([0-9]+)\\|([0-9A-Z]+)\\] (.*)$")
var onlineLocationRegexp = regexp.MustCompile("^(.*) \\(([0-9]+)\\|([0-9A-Z]+)\\)$")

type RWTHLocation struct {
	BuildingID *string
	RoomID     *string
	Name       string
}

func (l *RWTHLocation) Equal(o *RWTHLocation) bool {
	if l.BuildingID != nil && o.BuildingID != nil && l.RoomID != nil && o.RoomID != nil {
		// Don't match the name here, the building and room already give a unique identifier and names could be different (although the shouldn't)
		return *l.BuildingID == *o.BuildingID && *l.RoomID == *o.RoomID
	}

	return l.Name == o.Name
}

func (l *RWTHLocation) GetBuildingDetails() (*RWTHBuildingDetails, error) {
	if l.BuildingID != nil {
		return NewRWTHBuildingDetails(*l.BuildingID)
	}

	return nil, nil
}

func NewRWTHLocation(location string) *RWTHLocation {
	moodleLocationMatches := moodleLocationRegexp.FindStringSubmatch(location)
	if len(moodleLocationMatches) > 0 {
		return &RWTHLocation{
			BuildingID: stringP(moodleLocationMatches[1]),
			RoomID:     stringP(moodleLocationMatches[2]),
			Name:       moodleLocationMatches[3],
		}
	}

	onlineLocationMatches := onlineLocationRegexp.FindStringSubmatch(location)
	if len(onlineLocationMatches) > 0 {
		return &RWTHLocation{
			BuildingID: stringP(onlineLocationMatches[2]),
			RoomID:     stringP(onlineLocationMatches[3]),
			Name:       onlineLocationMatches[1],
		}
	}

	return &RWTHLocation{
		Name: location,
	}
}

func stringP(s string) *string {
	return &s
}
