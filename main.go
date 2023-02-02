package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/JenswBE/tuinfeestbeerse.be/preprocess"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Timezone = "Europe/Brussels"

func main() {
	// Setup logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load timezone
	timezone, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatal().Err(err).Str("timezone", Timezone).Msg("Failed to load timezone")
	}

	// Copy simple data files to destination
	srcPath := "./data"
	dstPath := "./website/data"
	preprocess.CopyDirContents(srcPath, dstPath)

	// Process timetable
	start, end := preprocess.GetEventStartAndEnd("./data/General.yml", timezone)
	srcTimetable := filepath.Join(srcPath, "Timetable.yml")
	dstTimetable := filepath.Join(dstPath, "Timetable.yml")
	preprocess.ProcessTimetable(srcTimetable, dstTimetable, start, end)
}
