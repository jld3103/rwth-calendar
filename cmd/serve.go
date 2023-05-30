package cmd

import (
	"fmt"
	ics "github.com/arran4/golang-ical"
	"github.com/provokateurin/rwth-calendar/internal"
	"github.com/spf13/cobra"
	"net/http"
	"regexp"
	"strings"
)

func NewServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				moodleValues, err := readMoodleValues(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				onlineValues, err := readOnlineValues(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				moodleResponse, err := http.Get(addQueryParametersToURL("https://moodle.rwth-aachen.de/calendar/export_execute.php", moodleValues))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer moodleResponse.Body.Close()

				moodleCalendar, err := ics.ParseCalendar(moodleResponse.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				onlineResponse, err := http.Get(addQueryParametersToURL("https://online.rwth-aachen.de/RWTHonlinej/ws/termin/ical", onlineValues))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer onlineResponse.Body.Close()

				onlineCalendar, err := ics.ParseCalendar(onlineResponse.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				mergedCalendar, err := internal.MergeCalendars(moodleCalendar, onlineCalendar)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = mergedCalendar.SerializeTo(w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Add("Content-Type", "text/calendar")
			})

			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				return fmt.Errorf("failed to run server: %w", err)
			}

			return nil
		},
	}

	return serveCmd
}

func addQueryParametersToURL(url string, values map[string]string) string {
	var parameters []string
	for key, value := range values {
		parameters = append(parameters, fmt.Sprintf("%s=%s", key, value))
	}

	return fmt.Sprintf("%s?%s", url, strings.Join(parameters, "&"))
}

var (
	numbersPattern              = regexp.MustCompile("^[0-9]+$")
	lowercaseLettersPattern     = regexp.MustCompile("^[a-z]+$")
	lowercaseHexadecimalPattern = regexp.MustCompile("^[a-f0-9]+$")
	uppercaseHexadecimalPattern = regexp.MustCompile("^[A-F0-9]+$")
)

func readMoodleValues(r *http.Request) (map[string]string, error) {
	return readValues(
		r,
		map[string]*regexp.Regexp{
			"userid":      numbersPattern,
			"authtoken":   lowercaseHexadecimalPattern,
			"preset_what": lowercaseLettersPattern,
			"preset_time": lowercaseLettersPattern,
		},
	)
}

func readOnlineValues(r *http.Request) (map[string]string, error) {
	return readValues(
		r,
		map[string]*regexp.Regexp{
			"pStud":  uppercaseHexadecimalPattern,
			"pToken": uppercaseHexadecimalPattern,
		},
	)
}

func readValues(r *http.Request, keysPatterns map[string]*regexp.Regexp) (map[string]string, error) {
	out := map[string]string{}
	for name, pattern := range keysPatterns {
		value := r.URL.Query().Get(name)
		if !pattern.MatchString(value) {
			return nil, fmt.Errorf("%q not valid", name)
		}

		out[name] = value
	}

	return out, nil
}
