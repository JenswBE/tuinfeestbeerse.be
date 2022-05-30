package generator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Data struct {
	CarouselImages map[string][]string
}

func GetData(staticDir string) Data {
	// Init data
	data := Data{CarouselImages: map[string][]string{}}

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
	return data
}
