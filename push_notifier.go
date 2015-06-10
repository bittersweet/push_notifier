package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ComputeHmac(message []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func handleReceive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("ReadAll", err)
	}

	computed := ComputeHmac(body, os.Getenv("HOOK_SECRET"))
	computedWithSha1 := fmt.Sprintf("sha1=%s", computed)
	signature := r.Header["X-Hub-Signature"][0]
	fmt.Println("computed: ", computed)
	fmt.Println("signature: ", signature)
	if computedWithSha1 != signature {
		log.Println("Computed message did not match signature")
		error := "Computed message did not match signature"
		http.Error(w, error, http.StatusUnauthorized)

		return
	}

	res := &Response{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("JSON unmarshal failed ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	message := res.formatMessage()
	fmt.Println("Sending message: ", message)
	files := res.getModifiedFiles()
	if len(files) > 0 {
		gF := res.getFileContents(files[0])
		fileBody := gF.decodedContent()
		sendMessage(message, fileBody)
	}

	responseObject := map[string]interface{}{
		"status": "ok",
	}

	output, err := json.MarshalIndent(responseObject, "", "  ")
	if err != nil {
		log.Fatal("Marshal Indent response ", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func getListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}

	return ":" + port, nil
}

func main() {
	http.HandleFunc("/receive", handleReceive)

	port, err := getListenAddress()
	if err != nil {
		log.Fatal("getListenAddress ", err)
	}

	fmt.Printf("Will start listening on http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}
