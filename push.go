package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type response struct {
	Ref        string
	Commits    []commit
	Repository repository
}

type author struct {
	Name string
}

type commit struct {
	ID      string
	Author  author
	Message string
	Added   []string
}

type repository struct {
	Name        string
	ContentsURL string `json:"contents_url"`
}

type githubFile struct {
	Content string
}

func (gf *githubFile) decodedContent() string {
	data, err := base64.StdEncoding.DecodeString(gf.Content)
	if err != nil {
		fmt.Println("base64 decoding error: ", err)
		return ""
	}
	return string(data)
}

func (r *response) formatMessage() string {
	// for _, commit := range r.Commits {
	// }
	if len(r.Commits) == 0 {
		return ""
	}

	commit := r.Commits[0]
	result := fmt.Sprintf("%s pushed to %s: %s", commit.Author.Name, r.Repository.Name, commit.Message)
	return result
}

func (r *response) getAddedFiles() []string {
	if len(r.Commits) == 0 {
		return []string{}
	}

	var mF []string
	for _, commit := range r.Commits {
		for _, file := range commit.Added {
			mF = append(mF, file)
		}
	}

	return mF
}

func (r *response) getFileContents(file string) *githubFile {
	re := regexp.MustCompile(`\{\+path\}`)
	fileURL := r.Repository.ContentsURL
	fileURL = re.ReplaceAllLiteralString(fileURL, file)
	token := os.Getenv("GITHUB_TOKEN")
	fileURL = fmt.Sprintf("%s?access_token=%s", fileURL, token)
	fmt.Println("Fetching GithubFile: ", fileURL)

	res, err := http.Get(fileURL)
	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	gF := &githubFile{}
	err = json.Unmarshal(contents, &gF)
	if err != nil {
		log.Println("JSON unmarshal failed ", err)
		return nil
	}

	return gF
}
