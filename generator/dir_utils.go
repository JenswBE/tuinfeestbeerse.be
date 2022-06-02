package generator

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func EnsureDirExists(dir string) {
	if len(dir) == 0 {
		log.Fatal().Msg("dir for EnsureDirExists cannot be an empty string")
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal().Err(err).Str("dir", dir).Msg("Failed to ensure directory exists")
	}
}

func DeleteDirContents(dir string) {
	if len(dir) == 0 {
		log.Fatal().Msg("dir for deleteDirContents cannot be an empty string")
	}
	mustRunBashCmd(fmt.Sprintf("rm -rf ./%s/*  ./%s/.* || true", dir, dir), "delete output dir contents")
}

func CopyDirContents(src, dst string) {
	mustRunBashCmd(fmt.Sprintf("cp -ar %s/* %s/", src, dst), "copy visible files from source to target directory")
	mustRunBashCmd(fmt.Sprintf("cp -ar %s/.??* %s/", src, dst), "copy hidden files from source to target directory")
}

func mustRunCmd(cmd *exec.Cmd, description string) {
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal().Err(err).Bytes("output", output).Msg("Failed to " + description)
	}
}

func mustRunBashCmd(bashCmd string, description string) {
	mustRunCmd(exec.Command("bash", "-c", bashCmd), description)
}
