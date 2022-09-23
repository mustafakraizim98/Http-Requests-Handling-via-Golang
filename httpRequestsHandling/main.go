package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

type JsonReplayHandler struct {
	Id        int    `json:"id"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Age       string `json:"age"`
	Email     string `json:"email"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func closeResponse(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func readReplyAsString(response *http.Response) string {
	_bytes, err := io.ReadAll(response.Body)
	check(err)

	var data strings.Builder
	_, err1 := data.Write(_bytes)
	check(err1)

	return data.String()
}

func readReplyByJsonKeys(response *http.Response) string {
	_bytes, err := io.ReadAll(response.Body)
	check(err)

	var jsonReplayHandler JsonReplayHandler
	if err := json.Unmarshal(_bytes, &jsonReplayHandler); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%+v", jsonReplayHandler)
}

func httpGetRequest(url string) {
	response, err := http.Get(url)
	check(err)

	defer closeResponse(response)
	replay := readReplyAsString(response)

	fmt.Println("StatusCode:", response.StatusCode)
	fmt.Println(replay)
}

func httpPostRequest(url string, data map[string]string) {
	jsonData, err := json.Marshal(data)
	check(err)

	response, err1 := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	check(err1)

	defer closeResponse(response)
	replay := readReplyByJsonKeys(response)

	fmt.Println("StatusCode:", response.StatusCode)
	fmt.Println(replay)
}

func main() {
	postId := "2"
	postIdParameter := "postId=" + url2.QueryEscape(postId)
	getUrl := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?%s", postIdParameter)
	httpGetRequest(getUrl)

	postUrl := "https://jsonplaceholder.typicode.com/posts"
	data := map[string]string{
		"id":         "1",
		"first-name": "Mustafa",
		"last-name":  "Kraizim",
		"age":        "23",
		"email":      "mustafakraizim@gmail.com",
	}
	httpPostRequest(postUrl, data)
}
