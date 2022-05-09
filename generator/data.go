package generator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Data struct {
	SliderImages map[string][]string
}

func GetData(staticDir string) Data {
	// Init data
	data := Data{SliderImages: map[string][]string{}}

	// Load slider images
	sliderImagesPath := filepath.Join(staticDir, "assets/images/sliders")
	filepath.WalkDir(sliderImagesPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		segments := strings.Split(path, string(os.PathSeparator))
		dirName := segments[len(segments)-2]
		data.SliderImages[dirName] = append(data.SliderImages[dirName], strings.TrimPrefix(path, staticDir+string(os.PathSeparator)))
		return nil
	})
	return data
}
