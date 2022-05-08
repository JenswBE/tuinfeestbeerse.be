package main

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func main() {
	deleteDirContents("output")
	copyDirContents("static", "output")
	parseTemplates("templates", "output")
}

func deleteDirContents(dir string) {
	if len(dir) == 0 {
		log.Fatal().Msg("dir for deleteDirContents cannot be an empty string")
	}
	mustRunBashCmd(fmt.Sprintf("rm -rf ./%s/*  ./%s/.* || true", dir, dir), "delete output dir contents")
}

func copyDirContents(src, dst string) {
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

func parseTemplates(templateDir, outputDir string) {
	// Init template and parse components
	templ := template.New("")
	templ.ParseGlob(path.Join(templateDir, "*.component.html"))

	// Parse and execute pages
	filepath.Walk(templateDir, func(templPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal().Err(err).Str("template", templPath).Msg("Failed to walk template directory")
		}
		if !strings.HasSuffix(templPath, ".page.html") {
			return nil
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
		templ.ExecuteTemplate(file, info.Name(), nil)
		return nil
	})
}
