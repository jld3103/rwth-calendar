package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arran4/golang-ical"
	"github.com/spf13/cobra"

	"github.com/jld3103/rwth-calendar/internal"
)

func NewGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:  "generate",
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs := args[0 : len(args)-1]
			output := args[len(args)-1]

			for _, path := range inputs {
				if !strings.HasSuffix(path, ".ics") {
					return fmt.Errorf("file has to be .ics: %s", path)
				}
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("file does not exist: %s", path)
				}
			}

			merged, err := parseCalendarFromPath(inputs[0])
			if err != nil {
				return err
			}

			for _, input := range inputs[1:] {
				var calendar *ics.Calendar
				calendar, err = parseCalendarFromPath(input)
				if err != nil {
					return err
				}
				merged, err = internal.MergeCalendars(merged, calendar)
				if err != nil {
					return err
				}
			}

			err = serializeCalendarToPath(merged, output)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return generateCmd
}

func parseCalendarFromPath(path string) (*ics.Calendar, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return nil, fmt.Errorf("failed to read ics file: %w", err)
	}

	calendar, err := ics.ParseCalendar(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ics file: %w", err)
	}

	return calendar, nil
}

func serializeCalendarToPath(calendar *ics.Calendar, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create ics file: %w", err)
	}

	err = calendar.SerializeTo(file)
	if err != nil {
		return fmt.Errorf("failed to serialize calendar: %w", err)
	}

	return nil
}
