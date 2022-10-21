package preprocess

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Timetable struct {
	// Show the timetable
	Show bool `yaml:"Show"`
	// Number of pixels to represent 1 minute on the timetable
	PixelsPerMinute float64 `yaml:"PixelsPerMinute"`
	// Number of minutes for each timeslot
	TimeslotMinutes int `yaml:"TimeslotMinutes"`
	// Locations at the event
	Locations []*TimetableLocation `yaml:"Locations"`

	// Computed - Don't set
	Slots []string `yaml:"Slots"`
	// Computed - Don't set
	SlotHeightPixels int `yaml:"SlotHeightPixels"`
}

type TimetableLocation struct {
	// Name of the location
	Name string `yaml:"Name"`
	// Shows at the location
	Shows []*TimetableShow `yaml:"Shows"`
}

type TimetableShow struct {
	// Name of the artist
	Name string `yaml:"Name"`
	// Start of the show. Use helper "calcTime".
	Start string `yaml:"Start"`
	// End of the show. Use helper "calcTime".
	End string `yaml:"End"`

	// Computed - Don't set
	StartPixels int `yaml:"StartPixels"`
	// Computed - Don't set
	HeightPixels int `yaml:"HeightPixels"`
}

func PreprocessTimetable(src, dst string, eventStart, eventEnd time.Time) {
	var timetable Timetable
	parseYAMLFile(src, &timetable)
	completeTimetable(&timetable, eventStart, eventEnd)
	writeYAMLFile(dst, timetable)
}

// completeTimetable completes and returns the timetable.
// Input is changed, but returned as well for convenience.
func completeTimetable(timetable *Timetable, eventStart, eventEnd time.Time) Timetable {
	// Calculate slots
	slots := make([]string, 0, int(eventEnd.Sub(eventEnd).Minutes())/30+1)
	for i := eventStart; i.Before(eventEnd) || i.Equal(eventEnd); i = i.Add(30 * time.Minute) {
		slots = append(slots, i.Format("15:04"))
	}
	timetable.Slots = slots
	timetable.SlotHeightPixels = int(float64(timetable.TimeslotMinutes) * timetable.PixelsPerMinute)

	// Calculate show starts and heights
	for _, loc := range timetable.Locations {
		for _, show := range loc.Shows {
			showStart := calcTime(eventStart, show.Start)
			if showStart.Before(eventStart) {
				log.Fatal().Time("event_start", eventStart).Time("show_start", showStart).Str("show_name", show.Name).Msg("Show before start of event")
			}
			showEnd := calcTime(eventStart, show.End)
			if showEnd.After(eventEnd) {
				log.Fatal().Time("event_end", eventStart).Time("show_end", showEnd).Str("show_name", show.Name).Msg("Show after end of event")
			}
			show.StartPixels = int(showStart.Sub(eventStart).Minutes()*timetable.PixelsPerMinute) + int(float64(timetable.TimeslotMinutes)/2.0)
			show.HeightPixels = int(showEnd.Sub(showStart).Minutes()*timetable.PixelsPerMinute) - 3 // Add some spacing between shows
		}
	}
	return *timetable
}

// calcTime calculates the time of the show based on the start
// date of the event and a simple time, e.g. 21:00.
// Times between 00:00 and 07:00 are considered the next day.
func calcTime(start time.Time, simpleTime string) time.Time {
	timeLog := log.Fatal().Str("time", simpleTime)
	timeParts := strings.Split(simpleTime, ":")
	if len(timeParts) != 2 {
		timeLog.Msg("Invalid simple time provided. Must be in format 00:00.")
	}
	hour, err := strconv.Atoi(strings.TrimPrefix(timeParts[0], "0"))
	if err != nil {
		timeLog.Str("hour", timeParts[0]).Msg("Invalid hour provided. Time be in format 00:00.")
	}
	mins, err := strconv.Atoi(strings.TrimPrefix(timeParts[1], "0"))
	if err != nil {
		timeLog.Str("minutes", timeParts[1]).Msg("Invalid minutes provided. Time be in format 00:00.")
	}
	timeDate := start
	if hour < 7 {
		timeDate = start.AddDate(0, 0, 1)
	}
	return time.Date(timeDate.Year(), timeDate.Month(), timeDate.Day(), hour, mins, 0, 0, start.Location())
}
