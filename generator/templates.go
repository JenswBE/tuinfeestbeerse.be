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
	// Log progress
	parseLog := log.With().Str("template_dir", templateDir).Str("output_dir", outputDir).Logger()
	parseLog.Info().Msg("Generating templates ...")
	parseLog.Debug().Interface("data", templateData).Msg("Generating templates with data...")

	// Init template and parse components
	rootTemplate := template.New("")
	rootTemplate.Funcs(template.FuncMap{
		"raw":           raw,
		"trimAfterDash": trimAfterDash,
	})
	rootTemplate.ParseGlob(path.Join(templateDir, "*component*.html"))

	// Parse and execute pages
	filepath.Walk(templateDir, func(templPath string, info os.FileInfo, err error) error {
		// Check if walk successful
		if err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to walk template directory")
		}

		// Strip template dir prefix
		if templPath == templateDir {
			return nil
		}
		templName := strings.TrimPrefix(templPath, templateDir)[1:] // Trim template dir and path separator

		// Create duplicate in output if directory
		walkLog := log.With().Str("template", templName).Logger()
		walkLog.Debug().Msg("Processing template ...")
		if info.IsDir() {
			walkLog.Debug().Msg("Template is a directory, duplicating in output directory")
			EnsureDirExists(filepath.Join(outputDir, templName))
			return nil
		}

		// Validate file type and name
		if !strings.HasSuffix(templPath, ".page.html") {
			return nil
		}
		templ, err := rootTemplate.Clone()
		if err != nil {
			walkLog.Fatal().Err(err).Msg("Failed to clone root template")
			return err
		}
		templContent, err := os.ReadFile(templPath)
		if err != nil {
			walkLog.Fatal().Err(err).Msg("Failed to read template file")
			return err
		}
		if _, err = templ.Parse(string(templContent)); err != nil {
			walkLog.Fatal().Err(err).Msg("Failed to parse template")
			return err
		}
		outputName := strings.TrimSuffix(templName, ".page.html") + ".html"
		file, err := os.Create(path.Join(outputDir, outputName))
		defer func() {
			file.Sync()
			if fileErr := file.Close(); fileErr != nil {
				walkLog.Fatal().Err(err).Msg("Failed to close parsed page output file")
			}
		}()
		if err != nil {
			walkLog.Fatal().Err(err).Msg("Failed to create parsed page output file")
			return err
		}
		walkLog.Info().Msg("Executing template ...")
		err = templ.Execute(file, templateData)
		if err != nil {
			walkLog.Fatal().Err(err).Msg("Failed to execute template")
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
