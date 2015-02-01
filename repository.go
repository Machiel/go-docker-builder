// repository.go
package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
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

	inputBuf, outputBuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)

	dir, err := os.Open(targetDir)
	defer dir.Close()

	if err != nil {
		log.Println("Unable to read dir: " + err.Error())
		return
	}

	files, err := dir.Readdir(0)

	if err != nil {
		log.Println("Unable to read dir files: " + err.Error())
		return
	}

	tr := tar.NewWriter(inputBuf)

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		fullPath := dir.Name() + string(filepath.Separator) + fileInfo.Name()

		log.Println("Packing: " + fullPath)

		file, err := os.Open(fullPath)

		if err != nil {
			log.Println("Something went wrong reading file " + fileInfo.Name())
			return
		}

		defer file.Close()

		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tr.WriteHeader(header)

		if err != nil {
			log.Println("Unable to write header " + err.Error())
			return
		}

		_, err = io.Copy(tr, file)

		if err != nil {
			log.Println("Unable to copy file to tar archive " + err.Error())
			return
		}

	}

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
