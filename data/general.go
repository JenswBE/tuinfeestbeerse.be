package data

import (
	"time"
)

const (
	TimeZone = "Europe/Brussels"
)

func getGeneral() General {
	return General{
		SaturdayStart: time.Date(2022, time.June, 25, 20, 0, 0, 0, getTimeZoneLocation()),
		SaturdayEnd:   time.Date(2022, time.June, 26, 3, 0, 0, 0, getTimeZoneLocation()),
		SundayStart:   time.Date(2022, time.June, 26, 13, 0, 0, 0, getTimeZoneLocation()),
		SundayEnd:     time.Date(2022, time.June, 26, 18, 0, 0, 0, getTimeZoneLocation()),
		Links: map[string]GeneralLink{
			"mail-info": {
				Title: "info@tuinfeestbeerse.be",
				Link:  "mailto:info@tuinfeestbeerse.be",
				Icon:  "info",
			},
			"mail-artists": {
				Title: "artiesten@tuinfeestbeerse.be",
				Link:  "mailto:artiesten@tuinfeestbeerse.be",
				Icon:  "music-note-beamed",
			},
			"facebook": {
				Title: "Facebook",
				Link:  "https://www.facebook.com/TuinfeestBeerse/",
				Icon:  "facebook",
			},
			"instagram": {
				Title: "Instagram",
				Link:  "https://www.instagram.com/tuinfeest_beerse/",
				Icon:  "instagram",
			},
			"youtube": {
				Title: "YouTube",
				Link:  "https://www.youtube.com/channel/UCdN0ff4rGCs9QGIcFIUlDXA",
				Icon:  "youtube",
			},
		},
	}
}

type General struct {
	SaturdayStart time.Time
	SaturdayEnd   time.Time
	SundayStart   time.Time
	SundayEnd     time.Time
	// List of link displayed at the bottom
	Links map[string]GeneralLink
}

type GeneralLink struct {
	Title string // Title of the link
	Icon  string // Name of icon. See https://icons.getbootstrap.com for supported names.
	Link  string // URL of the link
}
