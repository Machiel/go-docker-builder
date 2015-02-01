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

	_ = targetDir + "/output"

	log.Println("Running: docker build " + targetDir)
	buildCmd := exec.Command("docker", "build", targetDir)
	tpt, cerr := buildCmd.CombinedOutput()

	if cerr != nil {
		log.Println(cerr.Error())
		return
	}

	log.Println(string(tpt))

	log.Println("Build successful, output")

}
