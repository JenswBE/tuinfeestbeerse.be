package preprocess

import (
	"time"

	"github.com/rs/zerolog/log"
)

type General struct {
	SaturdayStart string `yaml:"SaturdayStart"`
	SaturdayEnd   string `yaml:"SaturdayEnd"`
}

func GetEventStartAndEnd(generalConfigSrc string, timezone *time.Location) (start time.Time, end time.Time) {
	// Parse general data
	var general General
	parseYAMLFile(generalConfigSrc, &general)

	// Parse dates
	start = parseGeneralDate("SaturdayStart", general.SaturdayStart, timezone)
	end = parseGeneralDate("SaturdayEnd", general.SaturdayEnd, timezone)
	return
}

func parseGeneralDate(name, input string, timezone *time.Location) time.Time {
	format := "2006-01-02T15:04:00" // RFC3339 without timezone
	value, err := time.ParseInLocation(format, input, timezone)
	if err != nil {
		log.Fatal().Err(err).Str("date_type", name).Str("input", input).Str("format", format).Msg("Failed to parse date")
	}
	return value
}
