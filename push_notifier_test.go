package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bittersweet/push_notifier/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("HOOK_SECRET", "secret")
}

func TestReceiveWithCorrectSignature(t *testing.T) {
	file, _ := ioutil.ReadFile("./fixtures/payload2.json")
	payload := strings.Split(string(file), "\n")[0] // strip off the newline
	url := "/receive"
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	request.Header.Add("X-Hub-Signature", "sha1=6f64f4be5185c80b52eea72fe120f9f0b5ce28c7")
	response := httptest.NewRecorder()
	handleReceive(response, request)

	assert.Equal(t, 200, response.Code)
}

func TestReceiveWithInCorrectSignature(t *testing.T) {
	payload := `{"payload": "incorrect"}`
	url := "/receive"
	request, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	request.Header.Add("X-Hub-Signature", "sha1=incorrectsignature")
	response := httptest.NewRecorder()
	handleReceive(response, request)

	assert.Equal(t, 401, response.Code)
}

func TestComputeHmac(t *testing.T) {
	body := []byte("body")
	result := computeHmac(body, "secret")

	expected := "a18991ff7e4513a1c2d2ee51e3a8e99ca891d9cd"
	assert.Equal(t, expected, result)
}
