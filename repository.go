// repository.go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// Repository holds data about the repository
type Repository struct {
	AbsoluteURL string `json:"absolute_url"`
	Fork        bool   `json:"fork"`
	IsPrivate   bool   `json:"is_private"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Scm         string `json:"scm"`
	Slug        string `json:"slug"`
	Website     string `json:"website"`
}

// GetRepoURL builds the URL of the repo
func (r Repository) GetRepoURL() string {
	var buffer bytes.Buffer

	buffer.WriteString("git@bitbucket.org:")
	buffer.WriteString(r.Owner)
	buffer.WriteString("/")
	buffer.WriteString(r.Slug)
	buffer.WriteString(".git")

	return buffer.String()
}

// StartBuild executes a build on the Commit Payload
func (r Repository) StartBuild() {

	cmd := exec.Command("git", "clone", r.GetRepoURL(), "clones/"+r.Slug)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %s", out.String())

}
