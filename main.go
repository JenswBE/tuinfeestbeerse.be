package main

import (
	"flag"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/JenswBE/go-pipeline/pipeline"
	"github.com/JenswBE/tuinfeestbeerse.be/preprocess"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Timezone = "Europe/Brussels"
	dataDir  = "data"
)

func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Setup logging
	logLevel := zerolog.InfoLevel
	if *debug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(consoleWriter)

	// Load timezone
	timezone, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatal().Err(err).Str("timezone", Timezone).Msg("Failed to load timezone")
	}

	// Process timetable
	start, end := preprocess.GetEventStartAndEnd(getDataPath("General.yml"), timezone)
	timetable := preprocess.ProcessTimetable(getDataPath("Timetable.yml"), start, end)

	// Copy static content
	cmd := exec.Command("cp", "-R", "static/.", "output/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal().Err(err).Str("cmd", cmd.String()).Bytes("output", output).Msg("Failed to copy static assets")
	}

	// Process website
	pipeline.
		NewHTML(template.FuncMap{
			"now":       time.Now,
			"parseTime": parseTimeInLocation(timezone),
			"pathJoin":  path.Join,
			"rawHTML":   rawHTML,
			"readDir":   readDir,
		}).
		WithTemplatesDir("templates").
		WithDataDir("data").
		WithOutputDir("output").
		LoadGlob("layout_*.gohtml").
		LoadGlob("page_index/component_*.gohtml").
		SetData("Timetable", timetable).
		SetDataYAML("Artists", "Artists.yml").
		SetDataYAML("General", "General.yml").
		SetDataYAML("Privacy", "Privacy.yml").
		SetDataYAML("Sponsors", "Sponsors.yml").
		LoadRenderSingle("page_404.gohtml", "404.html").
		LoadRenderSingle("page_index/page.gohtml", "index.html").
		LoadRenderSingle("page_huisreglement.gohtml", "huisreglement/index.html").
		LoadRenderSingle("page_vrijwilligers_privacy.gohtml", "vrijwilligers/privacy/index.html").
		Must()
}

func getDataPath(filename string) string {
	return filepath.Join(dataDir, filename)
}

func parseTimeInLocation(location *time.Location) func(string) time.Time {
	return func(input string) time.Time {
		parsed, err := time.ParseInLocation("2006-01-02T15:04:05", input, location)
		if err != nil {
			log.Fatal().Err(err).Str("input", input).Msg("Failed to parse given time string")
		}
		return parsed
	}
}

func rawHTML(input string) template.HTML {
	// #nosec G203
	return template.HTML(input)
}

func readDir(dir string) []fs.DirEntry {
	items, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal().Err(err).Str("dir", dir).Msg("Failed to read dir")
	}
	return items
}
