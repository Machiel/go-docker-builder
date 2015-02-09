// builder.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CommitPayload struct {
	Repository Repository `json:"repository"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Unable to read body " + err.Error())
		return
	}

	var commit CommitPayload
	err = json.Unmarshal(body, &commit)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid payload")
		return
	}

	go commit.Repository.StartBuild()
}

func main() {

	fmt.Println("Server started")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)

}
