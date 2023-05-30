package internal

import (
	"fmt"
	"github.com/arran4/golang-ical"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	moodleEvents := map[string][]*ics.VEvent{}
	for _, event := range events {
		cleanupProperty(event, ics.ComponentPropertySummary)
		cleanupProperty(event, ics.ComponentPropertyDescription)
		err := expandAllDayEvent(event)
		if err != nil {
			return nil, err
		}

		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		isMoodleStart := strings.HasSuffix(summary, " beginnt")
		isMoodleEnd := strings.HasSuffix(summary, " endet")
		if isMoodleStart || isMoodleEnd {
			splitSummary := strings.Split(summary, " ")
			name := strings.Join(splitSummary[:len(splitSummary)-1], " ")
			if val, ok := moodleEvents[name]; ok {
				moodleEvents[name] = append(val, event)
			} else {
				moodleEvents[name] = []*ics.VEvent{event}
			}
			continue
		}

		output.AddVEvent(event)
	}

	for name, moodleNamedEvents := range moodleEvents {
		if len(moodleNamedEvents) != 2 {
			for _, moodleNamedEvent := range moodleNamedEvents {
				output.AddVEvent(moodleNamedEvent)
			}

			continue
		}

		event, err := mergeMoodleNamedEvents(name, moodleNamedEvents[0], moodleNamedEvents[1])
		if err != nil {
			return nil, err
		}

		output.AddVEvent(event)
	}

	return output, nil
}

var escapeRegexp = regexp.MustCompile("\\\\([;,\\.])")

func cleanupProperty(event *ics.VEvent, property ics.ComponentProperty) {
	value := event.GetProperty(property).Value
	value = escapeRegexp.ReplaceAllString(value, "$1")
	value = html.UnescapeString(value)
	unquoted, err := strconv.Unquote(value)
	if err == nil {
		value = unquoted
	}

	event.SetProperty(property, value)
}

func mergeMoodleNamedEvents(name string, startEvent, endEvent *ics.VEvent) (*ics.VEvent, error) {
	event := startEvent
	event.SetSummary(name)

	allDayStart, err := startEvent.GetAllDayStartAt()
	if err != nil {
		var start time.Time
		start, err = startEvent.GetStartAt()
		if err != nil {
			return nil, err
		}
		event.SetStartAt(start)
	} else {
		event.SetAllDayStartAt(allDayStart)
	}

	allDayEnd, err := endEvent.GetAllDayEndAt()
	if err != nil {
		var end time.Time
		end, err = endEvent.GetEndAt()
		if err != nil {
			return nil, err
		}
		event.SetEndAt(end)
	} else {
		event.SetAllDayEndAt(allDayEnd)
	}

	return event, nil
}

func expandAllDayEvent(event *ics.VEvent) error {
	start, err := event.GetStartAt()
	if err != nil {
		return err
	}
	end, err := event.GetEndAt()
	if err != nil {
		return err
	}

	if start.Equal(end) {
		event.SetAllDayStartAt(start)
		event.SetAllDayEndAt(start.Add(time.Hour * 24))
	}

	return nil
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
