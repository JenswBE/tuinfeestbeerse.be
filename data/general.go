package data

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	TimeZone = "Europe/Brussels"
)

var SaturdayStart = time.Date(2022, time.June, 25, 20, 0, 0, 0, getTimeZoneLocation())
var SaturdayEnd = time.Date(2022, time.June, 26, 3, 0, 0, 0, getTimeZoneLocation())

type Data struct {
	CarouselImages map[string][]string
	Timetable      Timetable
}

func GetData(dataDir, staticDir string) (Data, error) {
	// Init data
	data := Data{
		CarouselImages: map[string][]string{},
	}

	// Load carousel images
	carouselImagesPath := filepath.Join(staticDir, "assets/images/carousels")
	filepath.WalkDir(carouselImagesPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		segments := strings.Split(path, string(os.PathSeparator))
		dirName := segments[len(segments)-2]
		data.CarouselImages[dirName] = append(data.CarouselImages[dirName], strings.TrimPrefix(path, staticDir+string(os.PathSeparator)))
		return nil
	})

	// Load other data
	data.Timetable = getTimetable(SaturdayStart, SaturdayEnd)
	return data, nil
}

func getTimeZoneLocation() *time.Location {
	location, err := time.LoadLocation(TimeZone)
	if err != nil {
		log.Fatal().Str("timezone", TimeZone).Msg("Failed to load timezone data")
	}
	return location
}
