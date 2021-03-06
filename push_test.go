package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bittersweet/push_notifier/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func pushPayload() *response {
	file, _ := ioutil.ReadFile("./fixtures/payload.json")

	res := &response{}
	err := json.Unmarshal(file, &res)
	if err != nil {
		log.Fatal("JSON unmarshal failed ", err)
	}
	return res
}

func TestPushFormatMessage(t *testing.T) {
	c := commit{
		ID:      "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
		Author:  author{"baxterthehacker"},
		Message: "Update README.md",
	}
	commits := make([]commit, 0, 1)
	commits = append(commits, c)
	repository := repository{
		Name: "public-repo",
	}
	r := &response{
		Ref:        "abcd",
		Commits:    commits,
		Repository: repository,
	}

	expected := "baxterthehacker pushed to public-repo: Update README.md"
	result := r.formatMessage()
	assert.Equal(t, expected, result)
}

func TestPushFormatMessageWithoutCommits(t *testing.T) {
	commits := make([]commit, 0, 1)
	r := &response{
		Ref:     "abcd",
		Commits: commits,
	}

	expected := ""
	result := r.formatMessage()
	assert.Equal(t, expected, result)
}

func TestAddedFiles(t *testing.T) {
	push := pushPayload()

	expectedFiles := []string{"README.md"}
	addedFiles := push.getAddedFiles()

	assert.Equal(t, len(addedFiles), 1)
	assert.Equal(t, addedFiles, expectedFiles)
}

func TestDecodedContent(t *testing.T) {
	gF := &githubFile{"IyBwdWJsaWMtcmVwbwo=\n"}

	expected := "# public-repo\n"
	result := gF.decodedContent()

	assert.Equal(t, expected, result)
}
