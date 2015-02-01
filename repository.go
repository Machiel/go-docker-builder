// repository.go
package main

import (
	"archive/tar"
	"bytes"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func getDockerClient() *docker.Client {
	endpoint := "tcp://192.168.59.103:2376"
	client, err := docker.NewClient(endpoint)

	if err != nil {
		log.Fatal("Unable to connect to docker: " + err.Error())
	}

	return client
}

type Repository struct {
	ID       int    `json:"id"`
	CloneURL string `json:"clone_url"`
}

// StartBuild executes a build on the Commit Payload
func (r Repository) StartBuild() {

	client := getDockerClient()

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

	t := time.Now()

	inputBuf, outputBuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)

	dockerFile, err := ioutil.ReadFile(targetDir + "/Dockerfile")

	if err != nil {
		log.Println("Unable to read dockerfile: " + err.Error())
		return
	}

	tr := tar.NewWriter(inputBuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(dockerFile)), ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write(dockerFile)
	tr.Close()

	opts := docker.BuildImageOptions{
		Name:         "made-by-db",
		InputStream:  inputBuf,
		OutputStream: outputBuf,
	}

	err = client.BuildImage(opts)

	if err != nil {
		log.Println("Unable to build: " + err.Error())
		return
	}

	log.Println(string(outputBuf.Bytes()))

}
