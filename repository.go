// repository.go
package main

import (
	"bytes"
	"fmt"
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

	cmd := exec.Command("git", "clone", r.CloneURL, "clones/"+repoID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Println(out.String())
		log.Fatal(err)
	}
	fmt.Printf("Result: %s", out.String())

}
