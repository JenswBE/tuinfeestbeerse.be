package main

import (
	"os"

	"github.com/JenswBE/tuinfeestbeerse.be/data"
	"github.com/JenswBE/tuinfeestbeerse.be/generator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Fetch data
	log.Info().Msg("Fetching data ...")
	data, err := data.GetData("data", "static")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get data")
	}

	// Refresh output dir
	log.Info().Msg("Refreshing output dir ...")
	generator.DeleteDirContents("output")
	generator.CopyDirContents("static", "output")

	// Generate templates
	log.Info().Interface("data", data).Msg("Generating templates ...")
	generator.ParseTemplates("templates", "output", data)
}
