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

type Response struct {
	Ref        string
	Commits    []Commit
	Repository Repository
}

type Author struct {
	Name string
}

type Commit struct {
	Id       string
	Author   Author
	Message  string
	Modified []string
}

type Repository struct {
	Name        string
	ContentsUrl string `json:"contents_url"`
}

type GithubFile struct {
	Content string
}

func (gf *GithubFile) decodedContent() string {
	data, err := base64.StdEncoding.DecodeString(gf.Content)
	if err != nil {
		fmt.Println("base64 decoding error: ", err)
		return ""
	}
	return string(data)
}

func (r *Response) formatMessage() string {
	// for _, commit := range r.Commits {
	// }
	if len(r.Commits) == 0 {
		return ""
	}

	commit := r.Commits[0]
	result := fmt.Sprintf("%s pushed to %s: %s", commit.Author.Name, r.Repository.Name, commit.Message)
	return result
}

func (r *Response) getModifiedFiles() []string {
	if len(r.Commits) == 0 {
		return []string{}
	}

	mF := make([]string, 0)
	for _, commit := range r.Commits {
		for _, file := range commit.Modified {
			mF = append(mF, file)
		}
	}

	return mF
}

func (r *Response) getFileContents(file string) *GithubFile {
	re := regexp.MustCompile(`\{\+path\}`)
	fileUrl := r.Repository.ContentsUrl
	fileUrl = re.ReplaceAllLiteralString(fileUrl, file)
	token := os.Getenv("GITHUB_TOKEN")
	fileUrl = fmt.Sprintf("%s?access_token=%s", fileUrl, token)
	fmt.Println("Fetching GithubFile: ", fileUrl)

	res, err := http.Get(fileUrl)
	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	gF := &GithubFile{}
	err = json.Unmarshal(contents, &gF)
	if err != nil {
		log.Println("JSON unmarshal failed ", err)
		return nil
	}

	return gF
}
