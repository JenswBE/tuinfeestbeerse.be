package main

import (
	"flag"
	"os"

	"github.com/JenswBE/tuinfeestbeerse.be/data"
	"github.com/JenswBE/tuinfeestbeerse.be/generator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "Sets log level to debug")
	flag.Parse()

	// Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Fetch data
	log.Info().Msg("Fetching data ...")
	data, err := data.GetData("data", "static")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get data")
	}

	// Refresh output dir
	log.Info().Msg("Refreshing output dir ...")
	generator.EnsureDirExists("output")
	generator.DeleteDirContents("output")
	generator.CopyDirContents("static", "output")

	// Generate templates
	generator.ParseTemplates("templates", "output", data)
}
