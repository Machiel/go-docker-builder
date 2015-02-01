// builder.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CommitPayload contains the data received by BitBucket
type CommitPayload struct {
	CanonURL string `json:"canon_url"`
	Commits  []struct {
		Author string `json:"author"`
		Branch string `json:"branch"`
		Files  []struct {
			File string `json:"file"`
			Type string `json:"type"`
		} `json:"files"`
		Message      string      `json:"message"`
		Node         string      `json:"node"`
		Parents      []string    `json:"parents"`
		RawAuthor    string      `json:"raw_author"`
		RawNode      string      `json:"raw_node"`
		Revision     interface{} `json:"revision"`
		Size         float64     `json:"size"`
		Timestamp    string      `json:"timestamp"`
		Utctimestamp string      `json:"utctimestamp"`
	} `json:"commits"`
	Repository Repository `json:"repository"`
	Truncated  bool       `json:"truncated"`
	User       string     `json:"user"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	payload := r.PostForm.Get("payload")

	var commit CommitPayload
	err := json.Unmarshal([]byte(payload), &commit)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid payload")
		return
	}

	go commit.Repository.StartBuild()

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	fmt.Println("Server started")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)

}
