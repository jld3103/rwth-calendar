package internal

import (
	"fmt"

	"github.com/arran4/golang-ical"
)

func MergeCalendars(calendar1, calendar2 *ics.Calendar) (*ics.Calendar, error) {
	output := ics.NewCalendar()
	output.SetMethod(ics.MethodPublish)
	output.SetVersion("2.0")
	output.SetProductId("rwth-calendar")

	events := calendar1.Events()
	for _, event2 := range calendar2.Events() {
		found := false
		for _, event1 := range events {
			same, err := areSameEvent(event2, event1)
			if err != nil {
				return nil, fmt.Errorf("failed to check if events are the same: %w", err)
			}

			if same {
				found = true
				break
			}
		}

		if !found {
			events = append(events, event2)
		}
	}

	for _, event := range events {
		output.AddVEvent(event)
	}

	return output, nil
}

func areSameEvent(event1, event2 *ics.VEvent) (bool, error) {
	sameStartAt, err := areSameStartAt(event1, event2)
	if err != nil {
		return false, err
	}

	sameEndAt, err := areSameEndAt(event1, event2)
	if err != nil {
		return false, err
	}

	sameLocation := areSameLocation(event1, event2)
	if err != nil {
		return false, err
	}

	return sameStartAt && sameEndAt && sameLocation, nil
}

func areSameStartAt(event1, event2 *ics.VEvent) (bool, error) {
	start1, err := event1.GetStartAt()
	if err != nil {
		return false, fmt.Errorf("failed to get start at: %w", err)
	}

	start2, err := event2.GetStartAt()
	if err != nil {
		return false, fmt.Errorf("failed to get start at: %w", err)
	}

	return start1.Equal(start2), nil
}

func areSameEndAt(event1, event2 *ics.VEvent) (bool, error) {
	end1, err := event1.GetEndAt()
	if err != nil {
		return false, fmt.Errorf("failed to get end at: %w", err)
	}

	end2, err := event2.GetEndAt()
	if err != nil {
		return false, fmt.Errorf("failed to get end at: %w", err)
	}

	return end1.Equal(end2), nil
}

func areSameLocation(event1, event2 *ics.VEvent) bool {
	var location1, location2 *RWTHLocation

	locationProperty1 := event1.GetProperty(ics.ComponentPropertyLocation)
	if locationProperty1 != nil {
		location1 = NewRWTHLocation(locationProperty1.Value)
	}

	locationProperty2 := event2.GetProperty(ics.ComponentPropertyLocation)
	if locationProperty2 != nil {
		location2 = NewRWTHLocation(locationProperty2.Value)
	}

	return location1 != nil && location2 != nil && location1.Equal(location2)
}
