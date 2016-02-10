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
	postURL := "https://slack.com/api/chat.postMessage"

	textTmpl := fmt.Sprintf("%s\n```\n%s\n```", message, fileContents)
	res, err := http.PostForm(postURL, url.Values{
		"token":   {os.Getenv("SLACK_TOKEN")},
		"channel": {os.Getenv("SLACK_CHANNEL")},
		"text":    {textTmpl},
	})
	defer res.Body.Close()

	slackBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Slack ReadAll error: ", err)
	}
	log.Println("Slack response: ", string(slackBody))
}
