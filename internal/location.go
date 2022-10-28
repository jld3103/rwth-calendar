package internal

import (
	"regexp"
)

var moodleLocationRegexp = regexp.MustCompile("^\\[([0-9]+)\\|([0-9A-Z]+)\\] (.*)$")
var onlineLocationRegexp = regexp.MustCompile("^(.*) \\(([0-9]+)\\|([0-9A-Z]+)\\)$")

type RWTHLocation struct {
	Building *string
	Room     *string
	Name     string
}

func (l *RWTHLocation) Equal(o *RWTHLocation) bool {
	if l.Building != nil && o.Building != nil && l.Room != nil && o.Room != nil {
		// Don't match the name here, the building and room already give a unique identifier and names could be different (although the shouldn't)
		return *l.Building == *o.Building && *l.Room == *o.Room
	}

	return l.Name == o.Name
}

func NewRWTHLocation(location string) *RWTHLocation {
	moodleLocationMatches := moodleLocationRegexp.FindStringSubmatch(location)
	if len(moodleLocationMatches) > 0 {
		return &RWTHLocation{
			Building: stringP(moodleLocationMatches[1]),
			Room:     stringP(moodleLocationMatches[2]),
			Name:     moodleLocationMatches[3],
		}
	}

	onlineLocationMatches := onlineLocationRegexp.FindStringSubmatch(location)
	if len(onlineLocationMatches) > 0 {
		return &RWTHLocation{
			Building: stringP(onlineLocationMatches[2]),
			Room:     stringP(onlineLocationMatches[3]),
			Name:     onlineLocationMatches[1],
		}
	}

	return &RWTHLocation{
		Name: location,
	}
}

func stringP(s string) *string {
	return &s
}
