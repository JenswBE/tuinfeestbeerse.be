package data

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Data struct {
	Artists        Artists
	CarouselImages map[string][]string
	SaturdayStart  time.Time
	SundayStart    time.Time
	General        General
	Sponsors       Sponsors
	Timetable      Timetable
}

func GetData(dataDir, staticDir string) (Data, error) {
	// Init data
	general := getGeneral()
	data := Data{
		Artists:        getArtists(),
		CarouselImages: map[string][]string{},
		General:        general,
		Sponsors:       getSponsors(),
		Timetable:      getTimetable(general.SaturdayStart, general.SaturdayEnd),
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
	return data, nil
}

func getTimeZoneLocation() *time.Location {
	location, err := time.LoadLocation(TimeZone)
	if err != nil {
		log.Fatal().Str("timezone", TimeZone).Msg("Failed to load timezone data")
	}
	return location
}
