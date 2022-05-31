package data

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	TimetablePixelsPerMinute float64 = 1.2 // Number of pixels to represent 1 minute on the timetable
	TimetableTimeslotMinutes float64 = 30  // Number of minutes for each timeslot
)

func getTimetable(eventStart, eventEnd time.Time) Timetable {
	return completeTimetable(&Timetable{
		Show: true,
		Locations: []*TimetableLocation{
			{
				Name: "Live on Stage",
				Shows: []*TimetableShow{
					{
						Name:  "Sub-lime",
						Start: calcTime(eventStart, "21:00"),
						End:   calcTime(eventStart, "22:15"),
					},
					{
						Name:  "The Skadillacs",
						Start: calcTime(eventStart, "22:45"),
						End:   calcTime(eventStart, "00:00"),
					},
					{
						Name:  "De Romeo's",
						Start: calcTime(eventStart, "00:30"),
						End:   calcTime(eventStart, "01:10"),
					},
					{
						Name:  "Koperdieven",
						Start: calcTime(eventStart, "01:10"),
						End:   calcTime(eventStart, "02:00"),
					},
				},
			},
			{
				Name: "Gazonfuif",
				Shows: []*TimetableShow{
					{
						Name:  "DJ Van het JH",
						Start: calcTime(eventStart, "20:00"),
						End:   calcTime(eventStart, "21:00"),
					},
					{
						Name:  "DJ N1co",
						Start: calcTime(eventStart, "21:00"),
						End:   calcTime(eventStart, "22:00"),
					},
					{
						Name:  "DJ Lenny",
						Start: calcTime(eventStart, "22:00"),
						End:   calcTime(eventStart, "23:30"),
					},
					{
						Name:  "Double-U",
						Start: calcTime(eventStart, "23:30"),
						End:   calcTime(eventStart, "01:00"),
					},
					{
						Name:  "DJ Polleke & Celis",
						Start: calcTime(eventStart, "01:00"),
						End:   calcTime(eventStart, "03:00"),
					},
				},
			},
			{
				Name: "Salsa",
				Shows: []*TimetableShow{
					{
						Name:  "Vzw Salsmanians",
						Start: calcTime(eventStart, "20:00"),
						End:   calcTime(eventStart, "02:00"),
					},
				},
			},
		},
	}, eventStart, eventEnd)
}

type Timetable struct {
	// Show the timetable
	Show bool
	// Locations at the event
	Locations []*TimetableLocation
	// Start of the event
	Start time.Time
	// end of the event
	End time.Time

	// Computed - Don't set
	Slots []string
	// Computed - Don't set
	SlotHeightPixels int
}

type TimetableLocation struct {
	// Name of the location
	Name string
	// Shows at the location
	Shows []*TimetableShow
}

type TimetableShow struct {
	// Name of the artist
	Name string
	// Start of the show. Use helper "calcTime".
	Start time.Time
	// End of the show. Use helper "calcTime".
	End time.Time

	// Computed - Don't set
	StartPixels int
	// Computed - Don't set
	HeightPixels int
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
	timetable.SlotHeightPixels = int(TimetableTimeslotMinutes * TimetablePixelsPerMinute)

	// Calculate show starts and heights
	for _, loc := range timetable.Locations {
		for _, show := range loc.Shows {
			if show.Start.Before(eventStart) {
				log.Fatal().Time("event_start", eventStart).Time("show_start", show.Start).Str("show_name", show.Name).Msg("Show before start of event")
			}
			if show.End.After(eventEnd) {
				log.Fatal().Time("event_end", eventStart).Time("show_end", show.Start).Str("show_name", show.Name).Msg("Show after end of event")
			}
			show.StartPixels = int(show.Start.Sub(eventStart).Minutes()*TimetablePixelsPerMinute) + int(TimetableTimeslotMinutes/2.0)
			show.HeightPixels = int(show.End.Sub(show.Start).Minutes()*TimetablePixelsPerMinute) - 3 // Add some spacing between shows
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
