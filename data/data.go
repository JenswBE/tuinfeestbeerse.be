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
	EventStart     time.Time
	General        General
	Sponsors       Sponsors
	Timetable      Timetable
}

func GetData(dataDir, staticDir string) (Data, error) {
	// Init data
	data := Data{
		Artists:        getArtists(),
		CarouselImages: map[string][]string{},
		EventStart:     SaturdayStart,
		General:        getGeneral(),
		Sponsors:       getSponsors(),
		Timetable:      getTimetable(SaturdayStart, SaturdayEnd),
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
