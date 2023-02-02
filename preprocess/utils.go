package preprocess

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

func CopyDirContents(src, dst string) {
	mustRunBashCmd(fmt.Sprintf("cp -ar %s/* %s/", src, dst), "copy visible files from source to target directory")
}

func mustRunCmd(cmd *exec.Cmd, description string) {
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal().Err(err).Bytes("output", output).Msg("Failed to " + description)
	}
}

func mustRunBashCmd(bashCmd string, description string) {
	mustRunCmd(exec.Command("bash", "-c", bashCmd), description)
}

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

func writeYAMLFile(dst string, data any) {
	logger := log.With().Str("destination_file", dst).Logger()
	file, err := os.Create(dst)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open destination YAML file")
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Fatal().Err(err).Msg("Failed to close destination YAML file")
		}
	}()

	// Write data into destination
	if err = yaml.NewEncoder(file).Encode(data); err != nil {
		logger.Fatal().Err(err).Msg("Failed to encode and write to YAML file")
	}
}
