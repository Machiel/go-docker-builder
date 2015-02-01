// repository.go
package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"time"
)

func NewClientFromEnv() (*docker.Client, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		return nil, fmt.Errorf("Missing DOCKER_HOST")
	}

	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""
	certPath := os.Getenv("DOCKER_CERT_PATH")

	if tlsVerify || certPath != "" {
		if certPath == "" {
			user, err := user.Current()
			if err != nil {
				return nil, err
			}

			certPath = path.Join(user.HomeDir, ".docker")
		}

		cert := path.Join(certPath, "cert.pem")
		key := path.Join(certPath, "key.pem")
		ca := ""
		if tlsVerify {
			ca = path.Join(certPath, "ca.pem")
		}

		return docker.NewTLSClient(endpoint, cert, key, ca)
	} else {
		return docker.NewClient(endpoint)
	}
}

type Repository struct {
	ID       int    `json:"id"`
	CloneURL string `json:"clone_url"`
}

// StartBuild executes a build on the Commit Payload
func (r Repository) StartBuild() {

	client, err := NewClientFromEnv()

	if err != nil {
		log.Println(err.Error())
		return
	}

	repoID := strconv.Itoa(r.ID)

	targetDir := "clones/" + repoID

	cmd := exec.Command("git", "clone", r.CloneURL, targetDir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

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

	log.Println("Successfully read Dockerfile")

	tr := tar.NewWriter(inputBuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(dockerFile)), ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write(dockerFile)
	tr.Close()

	opts := docker.BuildImageOptions{
		Name:         "made-by-db",
		InputStream:  inputBuf,
		OutputStream: outputBuf,
	}

	log.Println("Calling BuildImage")

	err = client.BuildImage(opts)

	log.Println("Finished calling BuildImage")

	if err != nil {
		log.Println("Unable to build: " + err.Error())
		return
	}

	log.Println("Survived without errors")

	log.Println(string(outputBuf.Bytes()))

}
