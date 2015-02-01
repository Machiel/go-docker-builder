// repository.go
package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
)

type Repository struct {
	ID       int    `json:"id"`
	CloneURL string `json:"clone_url"`
}

// StartBuild executes a build on the Commit Payload
func (r Repository) StartBuild() {

	repoID := strconv.Itoa(r.ID)

	targetDir := "clones/" + repoID

	cmd := exec.Command("git", "clone", r.CloneURL, targetDir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Println(err.Error())
		log.Println(out.String())
		return
	}

	log.Println("Clone successful")

	outputDir := targetDir + "/output"

	buildCmd := exec.Command("docker", "build", targetDir, "-v", outputDir+":/output")
	var buildCmdOut bytes.Buffer
	cmd.Stdout = &buildCmdOut
	err = buildCmd.Run()

	if err != nil {
		log.Println(err.Error())
		log.Println(buildCmdOut.String())
		return
	}

	log.Println("Build successful, output")
	log.Println(buildCmdOut.String())

}
