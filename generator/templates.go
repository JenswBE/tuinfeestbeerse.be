package generator

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/JenswBE/tuinfeestbeerse.be/data"
	"github.com/rs/zerolog/log"
)

func ParseTemplates(templateDir, outputDir string, templateData data.Data) {
	// Init template and parse components
	rootTemplate := template.New("")
	rootTemplate.Funcs(template.FuncMap{
		"raw":           raw,
		"trimAfterDash": trimAfterDash,
	})
	rootTemplate.ParseGlob(path.Join(templateDir, "*component*.html"))

	// Parse and execute pages
	filepath.Walk(templateDir, func(templPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to walk template directory")
		}
		if !strings.HasSuffix(templPath, ".page.html") {
			return nil
		}
		templ, err := rootTemplate.Clone()
		if err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to clone root template")
			return err
		}
		if _, err = templ.ParseFiles(templPath); err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to parse template")
			return err
		}
		outputName := strings.TrimSuffix(info.Name(), ".page.html") + ".html"
		file, err := os.Create(path.Join(outputDir, outputName))
		defer func() {
			file.Sync()
			if fileErr := file.Close(); fileErr != nil {
				log.Fatal().Err(err).Str("template", templPath).Msg("Failed to close parsed page output file")
			}
		}()
		if err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to create parsed page output file")
			return err
		}
		log.Info().Str("template", info.Name()).Msg("Executing template ...")
		err = templ.ExecuteTemplate(file, info.Name(), templateData)
		if err != nil {
			log.Fatal().Err(err).Str("template", info.Name()).Msg("Failed to execute template")
			return err
		}
		return nil
	})
}

func raw(input string) template.HTML {
	return template.HTML(input)
}

func trimAfterDash(input string) string {
	return strings.SplitN(input, "-", 1)[0]
}
