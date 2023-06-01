package preprocess

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

func parseYAMLFile(src string, target any) {
	// Read file
	yamlContent, err := os.ReadFile(src)
	if err != nil {
		log.Fatal().Err(err).Str("source_file", src).Msg("Failed to read source file contents")
	}

	// Parse into target
	if err = yaml.Unmarshal(yamlContent, target); err != nil {
		log.Fatal().Err(err).Str("source_file", src).Msg("Failed to unmarshal source file content")
	}
}
