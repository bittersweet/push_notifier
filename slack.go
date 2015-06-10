package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// https://api.slack.com/methods/chat.postMessage
func sendMessage(message string, fileContents string) {
	postUrl := "https://slack.com/api/chat.postMessage"

	textTmpl := fmt.Sprintf("%s\n```\n%s\n```", message, fileContents)
	res, err := http.PostForm(postUrl, url.Values{
		"token":   {os.Getenv("SLACK_TOKEN")},
		"channel": {os.Getenv("SLACK_CHANNEL")},
		"text":    {textTmpl},
	})
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ReadAll ", err)
	}
	// TODO: do something with result
}
