package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bittersweet/push_notifier/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestReceiveWithCorrectSignature(t *testing.T) {
	file, _ := ioutil.ReadFile("./fixtures/payload2.json")
	payload := strings.Split(string(file), "\n")[0] // strip off the newline
	url := "/receive"
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	request.Header.Add("X-Hub-Signature", "sha1=78dfb5f57f2673b1e006dddb21297834db260813")
	response := httptest.NewRecorder()
	handleReceive(response, request)

	assert.Equal(t, 200, response.Code)
}

func TestReceiveWithInCorrectSignature(t *testing.T) {
	payload := `{"payload": "incorrect"}`
	url := "/receive"
	request, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	request.Header.Add("X-Hub-Signature", "sha1=78dfb5f57f2673b1e006dddb21297834db260813")
	response := httptest.NewRecorder()
	handleReceive(response, request)

	assert.Equal(t, 401, response.Code)
}

func TestComputeHmac(t *testing.T) {
	body := []byte("body")
	result := ComputeHmac(body, "secret")

	expected := "a18991ff7e4513a1c2d2ee51e3a8e99ca891d9cd"
	assert.Equal(t, expected, result)
}